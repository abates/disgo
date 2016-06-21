package disgo

import (
	"errors"
)

type entry struct {
	locations []string
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

func (db *DB) AddFile(filename string) (phash PHash, err error) {
	phash, err = HashFile(filename)
	if err == nil {
		err = db.AddLocation(filename, phash)
	}
	return phash, err
}

func (db *DB) AddLocation(location string, phash PHash) error {
	e, found := db.entries[phash]
	if !found {
		e = new(entry)
		db.entries[phash] = e
	}
	e.locations = append(e.locations, location)
	return nil
}

func (db *DB) Find(phash PHash) (locations []string, err error) {
	if entry, found := db.entries[phash]; found {
		locations = entry.locations
	} else {
		err = ErrNotFound
	}

	return locations, err
}

func (db *DB) SearchByFile(filename string, maxDistance uint) (matches []string, err error) {
	h, err := HashFile(filename)
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
			results = append(results, e.locations...)
		}
	}

	return results, nil
}
