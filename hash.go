package disgo

import (
	"github.com/disintegration/imaging"
	"image"
)

type PHash uint64

func (p1 PHash) Distance(p2 PHash) (distance uint) {
	hamming := p1 ^ p2

	for hamming != 0 {
		distance += 1
		hamming &= hamming - 1
	}
	return
}

func intensity(img *image.NRGBA, row, column int) uint8 {
	offset := (row-img.Rect.Min.Y)*img.Stride + (column-img.Rect.Min.X)*4
	return uint8((uint16(img.Pix[offset]) + uint16(img.Pix[offset+1]) + uint16(img.Pix[offset+2])) / 3)
}

func HashFile(filename string) (PHash, error) {
	img, err := imaging.Open(filename)
	if err != nil {
		return 0, err
	}

	return Hash(img)
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
