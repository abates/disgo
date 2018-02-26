package disgo

import (
	"errors"
	"image"
	"io"
)

var (
	ErrNotFound         = errors.New("Image not found")
	ErrSaveNotSupported = errors.New("Underlying index does not support saving")
	ErrLoadNotSupported = errors.New("Underlying index does not support loading")
)

type Index interface {
	Insert(PHash) error
	Search(PHash, int) ([]PHash, error)
}

type Saveable interface {
	Save(io.Writer) error
}

type Loadable interface {
	Load(io.Reader) error
}

type DB struct {
	index       Index
	imageHasher func(image.Image) (PHash, error)
	fileHasher  func(io.Reader) (PHash, error)
}

func New() *DB {
	return NewDB(NewRadixIndex())
}

func NewDB(index Index) *DB {
	r := &DB{
		index:       index,
		imageHasher: Hash,
		fileHasher:  HashFile,
	}
	return r
}

func (db *DB) Add(img image.Image) (PHash, error) {
	hash, err := db.imageHasher(img)
	if err == nil {
		err = db.AddHash(hash)
	}
	return hash, err
}

func (db *DB) AddFile(reader io.Reader) (PHash, error) {
	hash, err := db.fileHasher(reader)
	if err == nil {
		err = db.AddHash(hash)
	}
	return hash, err
}

func (db *DB) AddHash(hash PHash) error {
	db.index.Insert(hash)
	return nil
}

func (db *DB) Save(writer io.Writer) error {
	if saver, ok := db.index.(Saveable); ok {
		return saver.Save(writer)
	}
	return ErrSaveNotSupported
}

func (db *DB) Load(reader io.Reader) error {
	if loader, ok := db.index.(Loadable); ok {
		return loader.Load(reader)
	}
	return ErrLoadNotSupported
}

func (db *DB) Search(img image.Image, maxDistance int) (matches []PHash, err error) {
	hash, err := db.imageHasher(img)
	if err == nil {
		matches, err = db.SearchByHash(hash, maxDistance)
	}
	return matches, err
}

func (db *DB) SearchByFile(reader io.Reader, maxDistance int) (matches []PHash, err error) {
	h, err := db.fileHasher(reader)
	if err == nil {
		matches, err = db.SearchByHash(h, maxDistance)
	}
	return matches, err
}

func (db *DB) SearchByHash(hash PHash, maxDistance int) ([]PHash, error) {
	return db.index.Search(hash, maxDistance)
}
