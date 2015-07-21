package disgo

import (
	"errors"
)

type entry struct {
	locations []string
}

type db struct {
	entries map[PHash]*entry
}

var (
	ErrNotFound = errors.New("Image not found")
)

var DB db

func init() {
	DB.entries = make(map[PHash]*entry)
}

func AddFile(filename string) error {
	phash, err := HashFile(filename)
	if err != nil {
		return err
	}
	return AddLocation(filename, phash)
}

func AddLocation(location string, phash PHash) error {
	e, found := DB.entries[phash]
	if !found {
		e = new(entry)
		DB.entries[phash] = e
	}
	e.locations = append(e.locations, location)
	return nil
}

func Find(phash PHash) ([]string, error) {
	if entry, found := DB.entries[phash]; found {
		return entry.locations, nil
	}
	return []string{}, ErrNotFound
}

func SearchByFile(filename string, maxDistance uint) ([]string, error) {
	h, e := HashFile(filename)
	if e != nil {
		return nil, e
	}

	return SearchByHash(h, maxDistance)
}

func SearchByHash(phash PHash, maxDistance uint) ([]string, error) {
	var results []string

	// look for existing entry within maxDistance of the hash
	for p, e := range DB.entries {
		if p.Distance(phash) <= maxDistance {
			results = append(results, e.locations...)
		}
	}

	return results, nil
}
