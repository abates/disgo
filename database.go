package disgo

import (
	"encoding"
	"errors"
	"image"
	"io"
	"io/ioutil"

	"github.com/disintegration/imaging"
)

var (
	ErrNotFound     = errors.New("Image not found")
	ErrNotSupported = errors.New("Underlying index does not support loading or saving")
)

type Index interface {
	Insert(PHash) error
	Search(PHash, int) ([]PHash, error)
}

type DB struct {
	index  Index
	hasher func(image.Image) (PHash, error)
}

func New() *DB {
	return NewDB(NewRadixIndex())
}

func NewDB(index Index) *DB {
	r := &DB{
		index:  index,
		hasher: Hash,
	}
	return r
}

func (db *DB) Add(img image.Image) (PHash, error) {
	hash, err := db.hasher(img)
	if err == nil {
		err = db.AddHash(hash)
	}
	return hash, err
}

func (db *DB) AddFile(reader io.Reader) (hash PHash, err error) {
	img, err := imaging.Decode(reader)
	if err == nil {
		hash, err = db.Add(img)
	}
	return hash, err
}

func (db *DB) AddHash(hash PHash) error {
	db.index.Insert(hash)
	return nil
}

func (db *DB) Save(writer io.Writer) error {
	buf, err := db.MarshalBinary()
	if err == nil {
		_, err = writer.Write(buf)
	}
	return err
}

func (db *DB) MarshalBinary() ([]byte, error) {
	if marshaler, ok := db.index.(encoding.BinaryMarshaler); ok {
		return marshaler.MarshalBinary()
	}
	return nil, ErrNotSupported
}

func (db *DB) Load(reader io.Reader) error {
	buf, err := ioutil.ReadAll(reader)
	if err == nil {
		err = db.UnmarshalBinary(buf)
	}
	return err
}

func (db *DB) UnmarshalBinary(buf []byte) error {
	if unmarshaler, ok := db.index.(encoding.BinaryUnmarshaler); ok {
		return unmarshaler.UnmarshalBinary(buf)
	}
	return ErrNotSupported
}

func (db *DB) Search(img image.Image, maxDistance int) (matches []PHash, err error) {
	hash, err := db.hasher(img)
	if err == nil {
		matches, err = db.SearchByHash(hash, maxDistance)
	}
	return matches, err
}

func (db *DB) SearchByFile(reader io.Reader, maxDistance int) (matches []PHash, err error) {
	img, err := imaging.Decode(reader)
	if err == nil {
		matches, err = db.Search(img, maxDistance)
	}
	return matches, err
}

func (db *DB) SearchByHash(hash PHash, maxDistance int) ([]PHash, error) {
	return db.index.Search(hash, maxDistance)
}
