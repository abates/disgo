package disgo

import (
	"errors"
	"io"
)

type entry struct {
	paths []string
}

type DB struct {
	entries map[PHash]*entry
}

var (
	ErrNotFound = errors.New("Image not found")
)

func NewDB() (db *DB) {
	db = new(DB)
	db.entries = make(map[PHash]*entry)
	return
}

func (db *DB) Add(path string, phash PHash) error {
	e, found := db.entries[phash]
	if !found {
		e = new(entry)
		db.entries[phash] = e
	}
	e.paths = append(e.paths, path)
	return nil
}

func (db *DB) Find(phash PHash) (paths []string, err error) {
	if entry, found := db.entries[phash]; found {
		paths = entry.paths
	} else {
		err = ErrNotFound
	}

	return paths, err
}

func (db *DB) SearchByFile(reader io.Reader, maxDistance uint) (matches []string, err error) {
	h, err := HashFile(reader)
	if err == nil {
		matches, err = db.SearchByHash(h, maxDistance)
	}
	return matches, err
}

func (db *DB) SearchByHash(phash PHash, maxDistance uint) ([]string, error) {
	var results []string

	// look for existing entry within maxDistance of the hash
	for p, e := range db.entries {
		if p.Distance(phash) <= maxDistance {
			results = append(results, e.paths...)
		}
	}

	return results, nil
}
