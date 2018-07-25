package disgo

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

type PHash uint64

func (p1 PHash) Distance(p2 PHash) (distance int) {
	hamming := p1 ^ p2

	for hamming != 0 {
		distance++
		hamming &= hamming - 1
	}
	return
}

func intensity(img image.Image, row, column int) uint8 {
	c := img.At(column, row)
	r, g, b, _ := c.RGBA()
	return uint8((r + g + b) / 3)
}

func Hash(img image.Image) (PHash, error) {
	rows := 8
	columns := 9
	var hash PHash

	img = imaging.Grayscale(img)
	img = imaging.Resize(img, columns, rows, imaging.Box)

	for row := 0; row < rows; row++ {
		for column := 0; column < columns-1; column++ {
			avg1 := intensity(img, row, column)
			avg2 := intensity(img, row, column+1)
			hash = hash << 1
			if avg1 > avg2 {
				hash = hash | 0x01
			}
		}
	}
	return hash, nil
}

func (h PHash) MarshalBinary() ([]byte, error) {
	return []byte{byte(h >> 56), byte(h >> 48), byte(h >> 40), byte(h >> 32), byte(h >> 24), byte(h >> 16), byte(h >> 8), byte(h)}, nil
}

func (h PHash) String() string {
	return fmt.Sprintf("0x%016x   %064b", uint64(h), uint64(h))
}
