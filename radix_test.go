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

type expectedNode struct {
	numChildren int
	length      uint8
	prefix      PHash
	path        []int
}

func nn(numChildren int, length uint8, prefix PHash, path ...int) expectedNode {
	return expectedNode{
		numChildren: numChildren,
		length:      length,
		prefix:      prefix,
		path:        path,
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		hash          PHash
		expectedNodes []expectedNode
	}{
		{0x4a00000000000000, []expectedNode{nn(0, 64, 0x4a00000000000000, 0)}},
		{0x5d00000000000000, []expectedNode{
			nn(2, 3, 0x4000000000000000, 0),
			nn(0, 61, 0x5000000000000000, 0, 0),
			nn(0, 61, 0xe800000000000000, 0, 1),
		}},
		{0x5900000000000000, []expectedNode{
			nn(2, 3, 0x4000000000000000, 0),
			nn(0, 61, 0x5000000000000000, 0, 0),
			nn(2, 2, 0xc000000000000000, 0, 1),
			nn(0, 59, 0x2000000000000000, 0, 1, 0),
			nn(0, 59, 0xa000000000000000, 0, 1, 1),
		}},
		{0x6900000000000000, []expectedNode{
			nn(2, 2, 0x4000000000000000, 0),
			nn(2, 1, 0x0000000000000000, 0, 0),
			nn(0, 61, 0x5000000000000000, 0, 0, 0),
			nn(2, 2, 0xc000000000000000, 0, 0, 1),
			nn(0, 59, 0x2000000000000000, 0, 0, 1, 0),
			nn(0, 59, 0xa000000000000000, 0, 0, 1, 1),
			nn(0, 62, 0xa400000000000000, 0, 1),
		}},
	}

	i := NewRadixIndex()
	for _, test := range tests {
		i.Insert(test.hash)
		for _, n := range test.expectedNodes {
			node := i.root
			for _, index := range n.path {
				if index == 0 {
					node = node.left
				} else {
					node = node.right
				}
			}

			numChildren := 0
			if node.left != nil {
				numChildren++
			}

			if node.right != nil {
				numChildren++
			}

			if n.numChildren != numChildren {
				t.Errorf("Expected %d children but got %d", n.numChildren, numChildren)
			}

			if n.length != node.length {
				t.Errorf("Expected node prefix length to be %d but it was %d", n.length, node.length)
			}

			if n.prefix != node.prefix {
				t.Errorf("Expected prefix to be 0x%016x but it was 0x%016x", n.prefix, node.prefix)
			}
		}
	}
}

func TestSearch(t *testing.T) {
	i := NewRadixIndex()
	i.Insert(0x4a00000000000000)
	i.Insert(0x5d00000000000000)
	i.Insert(0x5900000000000000)
	i.Insert(0x6900000000000000)

	hashes, _ := i.Search(0x4a00000000000000, 0)
	if len(hashes) != 1 {
		t.Errorf("Expected to retrieve an exact match")
	}
}
