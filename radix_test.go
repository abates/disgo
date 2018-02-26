package disgo

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNodeIsLeaf(t *testing.T) {
	tests := []struct {
		node     *Node
		expected bool
	}{
		{&Node{left: &Node{}}, false},
		{&Node{right: &Node{}}, false},
		{&Node{}, true},
	}

	for i, test := range tests {
		if test.node.IsLeaf() != test.expected {
			t.Errorf("tests[%d] expected %v got %v", i, test.expected, test.node.IsLeaf())
		}
	}
}

func TestNodeMatch(t *testing.T) {
	tests := []struct {
		p1          uint64
		l1          uint8
		p2          uint64
		expectedLen uint8
	}{
		{0xe900000000000000, 7, 0xef00000000000000, 5},
	}

	for i, test := range tests {
		node := &Node{prefix: PHash(test.p1), length: test.l1}
		length := node.Match(PHash(test.p2))
		if test.expectedLen != length {
			t.Errorf("tests[%d] expected %d got %d", i, test.expectedLen, length)
		}
	}
}

func TestNodeInsert(t *testing.T) {
	tests := []struct {
		startPrefix  uint64
		startLen     uint8
		insertPrefix uint64
		insertLen    uint8

		endPrefix      uint64
		endLen         uint8
		endLeftPrefix  uint64
		endLeftLen     uint8
		endRightPrefix uint64
		endRightLen    uint8
		resetPrevious  bool
	}{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, true},
		{0xff << 56, 64, 0xf0 << 56, 64, 0xf0 << 56, 4, 0x00, 60, 0x0f << 60, 60, true},
		{0x7f << 56, 64, 0xff << 56, 64, 0x00, 0, 0x7f << 56, 64, 0xff << 56, 64, true},

		// Initial child on the left, then build down the right
		{0x7f << 56, 4, 0x77f << 52, 12, 0x70 << 56, 4, 0x7f << 56, 8, 0, 0, true},
		{0, 0, 0x7f << 56, 8, 0x70 << 56, 4, 0x7f << 56, 8, 0x0f << 60, 4, false},
		{0, 0, 0x7ff << 52, 12, 0x70 << 56, 4, 0x7f << 56, 8, 0x0f << 60, 4, false},

		// Initial child on the right then build down the left
		{0xff << 56, 8, 0xfff << 52, 12, 0xff << 56, 8, 0x00, 0, 0xf << 60, 4, true},
		{0, 0, 0xff7 << 52, 16, 0xff << 56, 8, 0x7 << 60, 8, 0xf << 60, 4, false},
		{0, 0, 0xff7f << 48, 16, 0xff << 56, 8, 0x70 << 56, 4, 0xf << 60, 4, false},
	}

	root := &Node{}
	for i, test := range tests {
		if test.resetPrevious {
			root = &Node{prefix: PHash(test.startPrefix), length: test.startLen}
		}
		root.Insert(&Node{prefix: PHash(test.insertPrefix), length: test.insertLen})
		if root.length != test.endLen || root.prefix != PHash(test.endPrefix) {
			t.Errorf("tests[%d] expected %d:%x got %d:%x", i, test.endLen, test.endPrefix, root.length, uint64(root.prefix))
		}

		if root.left != nil {
			if root.left.length != test.endLeftLen || root.left.prefix != PHash(test.endLeftPrefix) {
				t.Errorf("tests[%d] expected %d:%x got %d:%x", i, test.endLeftLen, test.endLeftPrefix, root.left.length, uint64(root.left.prefix))
			}
		}

		if root.right != nil {
			if root.right.length != test.endRightLen || root.right.prefix != PHash(test.endRightPrefix) {
				t.Errorf("tests[%d] expected %d:%x got %d:%x", i, test.endRightLen, test.endRightPrefix, root.right.length, uint64(root.right.prefix))
			}
		}
	}
}

func TestNodeSearch(t *testing.T) {
	tests := []struct {
		prefixes        []uint64
		search          uint64
		match           uint64
		distance        int
		expectedMatches []PHash
	}{
		{nil, 0x00, 0x00, 0, nil},
		{[]uint64{0xff}, 0x00, 0x00, 0, nil},
		{[]uint64{0xff}, 0xff, 0x00, 0, []PHash{PHash(0xff)}},
	}

	for i, test := range tests {
		root := &Node{}
		for _, prefix := range test.prefixes {
			root.Insert(&Node{length: 64, prefix: PHash(prefix)})
		}

		matches := root.Search(PHash(test.search), PHash(test.match), test.distance)
		if !reflect.DeepEqual(test.expectedMatches, matches) {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedMatches, matches)
		}
	}
}

func TestNodeEncode(t *testing.T) {
	tests := []struct {
		prefixes []uint64
		expected []byte
	}{
		{[]uint64{0xff}, []byte{0x00, 0x40, 0, 0, 0, 0, 0, 0, 0, 0xff}},
		{[]uint64{0xfff << 52, 0xff7 << 52}, []byte{0xc0, 0x08, 0xff, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x38, 0x70, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x38, 0xf0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for i, test := range tests {
		root := &Node{}
		if len(test.prefixes) > 0 {
			root = &Node{length: 64, prefix: PHash(test.prefixes[0])}
			for _, prefix := range test.prefixes[1:] {
				root.Insert(&Node{length: 64, prefix: PHash(prefix)})
			}
		}
		writer := bytes.NewBuffer([]byte{})
		root.Encode(writer)
		buf := writer.Bytes()
		if !bytes.Equal(test.expected, buf) {
			t.Errorf("tests[%d] expected %x got %x", i, test.expected, buf)
		}
	}
}

func TestNodeEqual(t *testing.T) {
	tests := []struct {
		p1       []uint64
		p2       []uint64
		expected bool
	}{
		{nil, nil, true},
		{[]uint64{0}, nil, false},
		{[]uint64{0xff}, []uint64{0x00}, false},
		{[]uint64{0xff}, []uint64{0xff}, true},
		{[]uint64{0xff, 0xfe}, []uint64{0xff, 0xff}, false},
		{[]uint64{0xff, 0xfe}, []uint64{0xff, 0xfe}, true},
	}

	for i, test := range tests {
		var node1, node2 *Node
		if len(test.p1) > 0 {
			node1 = &Node{length: 64, prefix: PHash(test.p1[0])}
			for _, prefix := range test.p1[1:] {
				node1.Insert(&Node{length: 64, prefix: PHash(prefix)})
			}
		}

		if len(test.p2) > 0 {
			node2 = &Node{length: 64, prefix: PHash(test.p2[0])}
			for _, prefix := range test.p2[1:] {
				node2.Insert(&Node{length: 64, prefix: PHash(prefix)})
			}
		}

		if node1.Equal(node2) != test.expected {
			t.Errorf("tests[%d] expected %v got %v", i, test.expected, node1.Equal(node2))
		}
	}
}

func TestNodeDecode(t *testing.T) {
	tests := []struct {
		input    []byte
		prefixes []uint64
	}{
		{[]byte{0x00, 0x40, 0, 0, 0, 0, 0, 0, 0, 0xff}, []uint64{0xff}},
		{[]byte{0xc0, 0x08, 0xff, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x38, 0x70, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x38, 0xf0, 0, 0, 0, 0, 0, 0, 0}, []uint64{0xfff << 52, 0xff7 << 52}},
	}

	for i, test := range tests {
		expected := &Node{}
		if len(test.prefixes) > 0 {
			expected = &Node{length: 64, prefix: PHash(test.prefixes[0])}
			for _, prefix := range test.prefixes[1:] {
				expected.Insert(&Node{length: 64, prefix: PHash(prefix)})
			}
		}
		root := &Node{}
		reader := bytes.NewReader(test.input)
		root.Decode(reader)
		if !expected.Equal(root) {
			t.Errorf("tests[%d] expected %v got %v", i, expected, root)
		}
	}
}

type testRadixNode struct {
	insertedNode   *Node
	search         PHash
	searchMatch    PHash
	searchDistance int
	searchResults  []PHash
	encodeBuf      []byte
	decodeBuf      []byte
}

func (trn *testRadixNode) Insert(node *Node) { trn.insertedNode = node }

func (trn *testRadixNode) Search(search, match PHash, distance int) []PHash {
	trn.search = search
	trn.searchMatch = match
	trn.searchDistance = distance
	return trn.searchResults
}

func (trn *testRadixNode) Encode(writer io.Writer) error {
	_, err := writer.Write(trn.encodeBuf)
	return err
}

func (trn *testRadixNode) Decode(reader io.Reader) (err error) {
	trn.decodeBuf, err = ioutil.ReadAll(reader)
	return err
}

func TestRadixIndexInsert(t *testing.T) {
	index := NewRadixIndex()
	trn := &testRadixNode{}
	index.root = trn
	index.Insert(PHash(0x42))
	expected := &Node{length: 64, prefix: PHash(0x42)}
	if !expected.Equal(trn.insertedNode) {
		t.Errorf("Expected %v got %v", expected, trn.insertedNode)
	}
}

func TestRadixIndexSearch(t *testing.T) {
	index := NewRadixIndex()
	trn := &testRadixNode{}
	trn.searchResults = []PHash{PHash(0x42)}
	index.root = trn
	searchResults, _ := index.Search(PHash(0x53), 37)

	if trn.search != PHash(0x53) {
		t.Errorf("Expected %x got %x", PHash(0x53), trn.search)
	}

	if trn.searchMatch != 0x00 {
		t.Errorf("Expected 0x00 got %x", trn.searchMatch)
	}

	if trn.searchDistance != 37 {
		t.Errorf("Expected 37 got %d", trn.searchDistance)
	}

	if !reflect.DeepEqual(searchResults, trn.searchResults) {
		t.Errorf("Expected %v got %v", trn.searchResults, searchResults)
	}
}

func TestRadixSave(t *testing.T) {
	index := NewRadixIndex()
	trn := &testRadixNode{}
	trn.encodeBuf = []byte{0x73, 0x6f, 0x20, 0x6c, 0x6f, 0x6e, 0x67}
	index.root = trn

	buf := bytes.NewBuffer([]byte{})
	index.Save(buf)
	if !bytes.Equal(trn.encodeBuf, buf.Bytes()) {
		t.Errorf("Expected %s got %s", trn.encodeBuf, buf)
	}
}

func TestRadixIndexUnmarshal(t *testing.T) {
	index := NewRadixIndex()
	trn := &testRadixNode{}
	buf := []byte{0x61, 0x6e, 0x64, 0x20, 0x74, 0x68, 0x61, 0x6e, 0x6b, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x74, 0x68, 0x65, 0x20, 0x66, 0x69, 0x73, 0x68}
	reader := bytes.NewReader(buf)
	index.root = trn

	index.Load(reader)
	if !bytes.Equal(trn.decodeBuf, buf) {
		t.Errorf("Expected %s got %s", buf, trn.decodeBuf)
	}
}
