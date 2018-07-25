package disgo

import (
	"bytes"
	"fmt"
	"io"
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

type nodeFlags byte

func (flags nodeFlags) HasLeft() bool  { return flags&0x80 == 0x80 }
func (flags *nodeFlags) setHasLeft()   { (*flags) |= 0x80 }
func (flags nodeFlags) HasRight() bool { return flags&0x40 == 0x40 }
func (flags *nodeFlags) setHasRight()  { (*flags) |= 0x40 }

type radixNode interface {
	Insert(*Node)
	Search(PHash, PHash, int) []PHash
	Encode(io.Writer) error
	Decode(io.Reader) error
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
		n.prefix = n.prefix & bitmasks[matchLength]
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
	if n == nil {
		return nil
	}

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

	matches := n.left.Search(search, match, distance)
	matches = append(matches, n.right.Search(search, match, distance)...)
	return matches
}

func (n *Node) Encode(writer io.Writer) error {
	if n == nil {
		return nil
	}

	buf := make([]byte, 10)
	flags := nodeFlags(0x00)
	if n.left != nil {
		flags.setHasLeft()
	}

	if n.right != nil {
		flags.setHasRight()
	}

	buf[0] = byte(flags)
	buf[1] = byte(n.length)
	buf[2] = byte(n.prefix >> 56)
	buf[3] = byte(n.prefix >> 48)
	buf[4] = byte(n.prefix >> 40)
	buf[5] = byte(n.prefix >> 32)
	buf[6] = byte(n.prefix >> 24)
	buf[7] = byte(n.prefix >> 16)
	buf[8] = byte(n.prefix >> 8)
	buf[9] = byte(n.prefix)

	_, err := writer.Write(buf)
	if err == nil {
		err = n.left.Encode(writer)
	}

	if err == nil {
		err = n.right.Encode(writer)
	}
	return err
}

func (n *Node) Decode(reader io.Reader) error {
	buf := make([]byte, 10)
	_, err := reader.Read(buf)
	if err == nil {
		n.left = nil
		n.right = nil

		flags := nodeFlags(buf[0])
		n.length = buf[1]
		n.prefix = PHash(buf[2]) << 56
		n.prefix |= PHash(buf[3]) << 48
		n.prefix |= PHash(buf[4]) << 40
		n.prefix |= PHash(buf[5]) << 32
		n.prefix |= PHash(buf[6]) << 24
		n.prefix |= PHash(buf[7]) << 16
		n.prefix |= PHash(buf[8]) << 8
		n.prefix |= PHash(buf[9])

		if flags.HasLeft() {
			n.left = &Node{}
			n.left.Decode(reader)
		}

		if flags.HasRight() {
			n.right = &Node{}
			n.right.Decode(reader)
		}
	}
	return err
}

func (n *Node) Equal(other *Node) bool {
	if n == other {
		return true
	} else if n == nil || other == nil {
		return false
	}

	if n.prefix == other.prefix && n.length == other.length {
		return n.left.Equal(other.left) && n.right.Equal(other.right)
	}
	return false
}

func (n *Node) String() string {
	return fmt.Sprintf("%2d   %v", n.length, n.prefix)
}

type RadixIndex struct {
	root radixNode
}

func NewRadixIndex() *RadixIndex {
	ri := new(RadixIndex)
	ri.root = &Node{}
	return ri
}

func (ri *RadixIndex) Insert(hash PHash) error {
	node := &Node{
		prefix: hash,
		length: 64,
	}
	ri.root.Insert(node)
	return nil
}

func (ri *RadixIndex) Search(hash PHash, distance int) ([]PHash, error) {
	return ri.root.Search(hash, 0x00, distance), nil
}

func (ri *RadixIndex) MarshalBinary() ([]byte, error) {
	writer := bytes.NewBuffer(nil)
	err := ri.root.Encode(writer)
	return writer.Bytes(), err
}

func (ri *RadixIndex) UnmarshalBinary(buf []byte) error {
	reader := bytes.NewReader(buf)
	return ri.root.Decode(reader)
}
