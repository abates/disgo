package disgo

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io"
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

func intensity(img *image.NRGBA, row, column int) uint8 {
	offset := (row-img.Rect.Min.Y)*img.Stride + (column-img.Rect.Min.X)*4
	return uint8((uint16(img.Pix[offset]) + uint16(img.Pix[offset+1]) + uint16(img.Pix[offset+2])) / 3)
}

func HashFile(reader io.Reader) (hash PHash, err error) {
	img, err := imaging.Decode(reader)
	if err == nil {
		hash, err = Hash(img)
	}

	return hash, err
}

func Hash(img image.Image) (PHash, error) {
	rows := 8
	columns := 9
	var hash PHash

	grayscale := imaging.Grayscale(img)
	grayscale = imaging.Resize(grayscale, columns, rows, imaging.Box)

	for row := 0; row < rows; row++ {
		for column := 0; column < columns-1; column++ {
			avg1 := intensity(grayscale, row, column)
			avg2 := intensity(grayscale, row, column+1)
			hash = hash << 1
			if avg1 > avg2 {
				hash = hash | 0x01
			}
		}
	}
	return hash, nil
}

func (h PHash) String() string {
	return fmt.Sprintf("0x%016x   %064b", uint64(h), uint64(h))
}
