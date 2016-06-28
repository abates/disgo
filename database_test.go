package disgo

import (
	"os"
	"sort"
	"testing"
)

func addFile(db *DB, path string) PHash {
	file, _ := os.Open(path)
	phash, _ := HashFile(file)
	db.Add(path, phash)
	return phash
}

func TestFind(t *testing.T) {
	db := NewDB()
	hash1 := addFile(db, "images/gopher1.png")
	hash2 := addFile(db, "images/gopher2.png")

	entries, _ := db.Find(hash1)
	if len(entries) != 2 {
		t.Logf("Expected to find two entries for hash %v but only got %d", hash1, len(entries))
		t.Fail()
	}

	entries, _ = db.Find(hash2)
	if len(entries) != 2 {
		t.Logf("Expected to find two entries for hash %v but only got %d", hash2, len(entries))
		t.Fail()
	}

	file, _ := os.Open("images/ascendingGradient.png")
	hash3, _ := HashFile(file)
	entries, err := db.Find(hash3)
	if len(entries) > 0 {
		t.Logf("Expected not to find hash %v but found %d entries", hash3, len(entries))
		t.Fail()
	}

	if err != ErrNotFound {
		t.Logf("Expected %v but got %v", ErrNotFound, err)
		t.Fail()
	}
}

func TestSearchByFile(t *testing.T) {
	db := NewDB()
	addFile(db, "images/ascendingGradient.png")
	addFile(db, "images/descendingGradient.png")
	addFile(db, "images/alternatingGradient.png")
	addFile(db, "images/gopher1.png")
	addFile(db, "images/gopher2.png")

	file, _ := os.Open("images/gopher3.png")
	paths, err := db.SearchByFile(file, 5)
	if err != nil {
		t.Logf("Expected no error while searching by file.  Got: %v", err)
		t.Fail()
	}

	if len(paths) != 2 {
		t.Logf("Expected exactly two matching images.  Got: %d", len(paths))
		t.FailNow()
	}

	sort.Strings(paths)
	if paths[0] != "images/gopher1.png" {
		t.Logf("Expected to match images/gopher1.png but got %s instead", paths[0])
		t.Fail()
	}

	if paths[1] != "images/gopher2.png" {
		t.Logf("Expected to match images/gopher2.png but got %s instead", paths[1])
		t.Fail()
	}
}
