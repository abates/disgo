package disgo

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMatch(t *testing.T) {
	n := &Node{
		prefix: 0xe900000000000000,
		length: 7,
	}

	v := PHash(0xef00000000000000)
	length := n.Match(v)

	if length != 5 {
		t.Errorf("Expected length 5 got %d", length)
	}
}

func up(frame int, s string) string {
	_, origFile, origLine, _ := runtime.Caller(1)
	_, frameFile, frameLine, _ := runtime.Caller(frame + 1)
	if origFile != frameFile {
		return s // Deferred call after a panic or runtime.Goexit()
	}
	erase := []byte("\b\b\b")
	for ; origLine > 9; origLine /= 10 {
		erase = append(erase, '\b')
	}
	return fmt.Sprintf("%s%d: %s", erase, frameLine, s)
}

func checkValue(t *testing.T, n *Node, expectedLen uint8, expectedValue PHash) {
	if n.prefix != expectedValue {
		t.Errorf(up(1, "Expected 0x%016x, but got 0x%016x"), expectedValue, n.prefix)
	}

	if n.length != expectedLen {
		t.Errorf(up(1, "Expected length %d but got %d"), expectedLen, n.length)
	}
}

func TestInsert(t *testing.T) {
	i := NewIndex()
	i.Insert(0x4a00000000000000)

	if len(i.root.children) < 1 {
		t.Errorf("Expected the new value to be inserted into the root node")
	}

	checkValue(t, i.root.children[0], 64, 0x4a00000000000000)

	i.Insert(0x5d00000000000000)

	checkValue(t, i.root.children[0], 3, 0x4000000000000000)
	checkValue(t, i.root.children[0].children[0], 61, 0x5000000000000000)
	checkValue(t, i.root.children[0].children[1], 61, 0xe800000000000000)

	i.Insert(0x5900000000000000)

	checkValue(t, i.root.children[0], 3, 0x4000000000000000)
	checkValue(t, i.root.children[0].children[0], 61, 0x5000000000000000)
	checkValue(t, i.root.children[0].children[1], 2, 0xc000000000000000)
	checkValue(t, i.root.children[0].children[1].children[0], 59, 0xa000000000000000)
	checkValue(t, i.root.children[0].children[1].children[1], 59, 0x2000000000000000)

	i.Insert(0x6900000000000000)
	checkValue(t, i.root.children[0], 2, 0x4000000000000000)
	checkValue(t, i.root.children[0].children[0], 1, 0x0000000000000000)
	checkValue(t, i.root.children[0].children[0].children[0], 61, 0x5000000000000000)
	checkValue(t, i.root.children[0].children[0].children[1], 2, 0xc000000000000000)
	checkValue(t, i.root.children[0].children[0].children[1].children[0], 59, 0xa000000000000000)
	checkValue(t, i.root.children[0].children[0].children[1].children[1], 59, 0x2000000000000000)

	checkValue(t, i.root.children[0].children[1], 62, 0xa400000000000000)
}

func TestSearch(t *testing.T) {
	i := NewIndex()
	i.Insert(0x4a00000000000000)
	i.Insert(0x5d00000000000000)
	i.Insert(0x5900000000000000)
	i.Insert(0x6900000000000000)

	hashes, _ := i.Search(0x4a00000000000000, 0)
	if len(hashes) != 1 {
		t.Errorf("Expected to retrieve an exact match")
	}
}
