package disgo

import (
	"errors"
	"io"
)

var (
	ErrNotFound = errors.New("Image not found")
)

type Index interface {
	Insert(PHash) error
	Search(PHash, int) ([]PHash, error)
}

type DB struct {
	index Index
	paths map[PHash][]string
}

func New() *DB {
	return NewDB(NewRadixIndex())
}

func NewDB(index Index) *DB {
	r := new(DB)
	r.index = index
	r.paths = make(map[PHash][]string)
	return r
}

func (db *DB) Add(path string, hash PHash) error {
	if db.paths[hash] == nil {
		db.paths[hash] = []string{path}
	} else {
		db.paths[hash] = append(db.paths[hash], path)
	}
	db.index.Insert(hash)
	return nil
}

func (db *DB) Find(hash PHash) (paths []string, err error) {
	if db.paths[hash] == nil {
		return nil, ErrNotFound
	}
	return db.paths[hash], nil
}

func (db *DB) SearchByFile(reader io.Reader, maxDistance int) (matches []string, err error) {
	h, err := HashFile(reader)
	if err == nil {
		matches, err = db.Search(h, maxDistance)
	}
	return matches, err
}

func (db *DB) Search(hash PHash, maxDistance int) ([]string, error) {
	var results []string
	//fmt.Printf("Search:      %v\n", hash)
	hashes, _ := db.index.Search(hash, maxDistance)

	for _, hash := range hashes {
		results = append(results, db.paths[hash]...)
	}
	return results, nil
}
