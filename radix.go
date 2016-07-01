package disgo

import (
	"fmt"
)

var bitmasks = []PHash{
	0x0000000000000000, 0x8000000000000000, 0xc000000000000000, 0xe000000000000000, 0xf000000000000000,
	0xf800000000000000, 0xfc00000000000000, 0xfe00000000000000, 0xff00000000000000, 0xff80000000000000,
	0xffc0000000000000, 0xffe0000000000000, 0xfff0000000000000, 0xfff8000000000000, 0xfffc000000000000,
	0xfffe000000000000, 0xffff000000000000, 0xffff800000000000, 0xffffc00000000000, 0xffffe00000000000,
	0xfffff00000000000, 0xfffff80000000000, 0xfffffc0000000000, 0xfffffe0000000000, 0xffffff0000000000,
	0xffffff8000000000, 0xffffffc000000000, 0xffffffe000000000, 0xfffffff000000000, 0xfffffff800000000,
	0xfffffffc00000000, 0xfffffffe00000000, 0xffffffff00000000, 0xffffffff80000000, 0xffffffffc0000000,
	0xffffffffe0000000, 0xfffffffff0000000, 0xfffffffff8000000, 0xfffffffffc000000, 0xfffffffffe000000,
	0xffffffffff000000, 0xffffffffff800000, 0xffffffffffc00000, 0xffffffffffe00000, 0xfffffffffff00000,
	0xfffffffffff80000, 0xfffffffffffc0000, 0xfffffffffffe0000, 0xffffffffffff0000, 0xffffffffffff8000,
	0xffffffffffffc000, 0xffffffffffffe000, 0xfffffffffffff000, 0xfffffffffffff800, 0xfffffffffffffc00,
	0xfffffffffffffe00, 0xffffffffffffff00, 0xffffffffffffff80, 0xffffffffffffffc0, 0xffffffffffffffe0,
	0xfffffffffffffff0, 0xfffffffffffffff8, 0xfffffffffffffffc, 0xfffffffffffffffe, 0xffffffffffffffff,
}

type Node struct {
	prefix PHash
	length uint8
	left   *Node
	right  *Node
}

func (n *Node) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *Node) String() string {
	return fmt.Sprintf("%2d   %v", n.length, n.prefix)
}

func (n *Node) Insert(value *Node) {
	if n.length == value.length && n.prefix == value.prefix {
		return
	}

	matchLength := n.Match(value.prefix)
	value.prefix = value.prefix << matchLength
	value.length = value.length - matchLength

	if n.length > matchLength {
		newNode := &Node{
			prefix: n.prefix << matchLength,
			length: n.length - matchLength,
			left:   n.left,
			right:  n.right,
		}

		n.length = matchLength
		n.prefix = n.prefix & bitmasks[matchLength]
		if value.prefix&bitmasks[1] == 0 {
			n.left = value
			n.right = newNode
		} else {
			n.left = newNode
			n.right = value
		}
	} else if n.IsLeaf() {
		if value.prefix&bitmasks[1] == 0 {
			n.left = value
		} else {
			n.right = value
		}
	} else {
		if value.prefix&bitmasks[1] == 0 {
			if n.left == nil {
				n.left = value
			} else {
				n.left.Insert(value)
			}
		} else {
			if n.right == nil {
				n.right = value
			} else {
				n.right.Insert(value)
			}
		}
	}
}

func (n *Node) distance(hash PHash) int {
	return n.prefix.Distance(hash & bitmasks[n.length])
}

func (n *Node) Search(search PHash, match PHash, distance int) []PHash {
	// add my prefix to the match
	if n.length > 0 {
		match = match << n.length
		match = match | n.prefix>>(64-n.length)
		distance -= n.distance(search)
		search = search << n.length
	}

	if distance < 0 {
		return nil
	}

	if n.length > 0 && n.IsLeaf() {
		return []PHash{match}
	}

	var leftMatches, rightMatches []PHash
	if n.left != nil {
		leftMatches = n.left.Search(search, match, distance)
	}

	if n.right != nil {
		rightMatches = n.right.Search(search, match, distance)
	}

	if leftMatches != nil && rightMatches != nil {
		return append(leftMatches, rightMatches...)
	} else if leftMatches != nil {
		return leftMatches
	}
	return rightMatches
}

func (n *Node) Match(value PHash) (length uint8) {
	length = n.length
	for length = n.length; length > 0; length-- {
		mask := bitmasks[length]
		if value&mask == n.prefix&mask {
			break
		}
	}
	return
}

func (n *Node) Lookup(value PHash) (*Node, uint8) {
	length := n.Match(value)
	value = value << length
	if length < n.length {
		return n, length
	} else if value&bitmasks[1] == 0 {
		if n.left == nil {
			return n, length
		}
		match, l := n.left.Lookup(value)
		return match, length + l
	}
	if n.right == nil {
		return n, length
	}
	match, l := n.right.Lookup(value)
	return match, length + l
}

type RadixIndex struct {
	root *Node
}

func NewRadixIndex() *RadixIndex {
	i := new(RadixIndex)
	i.root = new(Node)
	i.root.prefix = 0
	i.root.length = 0
	return i
}

func (i *RadixIndex) Insert(hash PHash) error {
	node := &Node{
		prefix: hash,
		length: 64,
	}
	i.root.Insert(node)
	return nil
}

func (i *RadixIndex) Search(hash PHash, distance int) ([]PHash, error) {
	return i.root.Search(hash, 0x00, distance), nil
}
