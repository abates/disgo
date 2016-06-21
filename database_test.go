package disgo

import (
	"sort"
	"testing"
)

func TestAddFile(t *testing.T) {
	db := NewDB()

	_, err := db.AddFile("images/ascendingGradient.png")
	if err != nil {
		t.Logf("Expected no error while adding file.  Got: %v", err)
		t.Fail()
	}

	if len(db.entries) != 1 {
		t.Log("Database should contain exactly one entry")
		t.Fail()
	}

}

func TestFind(t *testing.T) {
	db := NewDB()
	hash1, _ := db.AddFile("images/gopher1.png")
	hash2, _ := db.AddFile("images/gopher2.png")

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

	hash3, _ := HashFile("images/ascendingGradient.png")
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
	db.AddFile("images/ascendingGradient.png")
	db.AddFile("images/descendingGradient.png")
	db.AddFile("images/alternatingGradient.png")
	db.AddFile("images/gopher1.png")
	db.AddFile("images/gopher2.png")

	filenames, err := db.SearchByFile("images/gopher3.png", 5)
	if err != nil {
		t.Logf("Expected no error while searching by file.  Got: %v", err)
		t.Fail()
	}

	if len(filenames) != 2 {
		t.Logf("Expected exactly two matching images.  Got: %d", len(filenames))
		t.FailNow()
	}

	sort.Strings(filenames)
	if filenames[0] != "images/gopher1.png" {
		t.Logf("Expected to match images/gopher1.png but got %s instead", filenames[0])
		t.Fail()
	}

	if filenames[1] != "images/gopher2.png" {
		t.Logf("Expected to match images/gopher2.png but got %s instead", filenames[1])
		t.Fail()
	}
}
