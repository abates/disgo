package disgo

import (
	"sort"
	"testing"
)

func TestAddFile(t *testing.T) {
	if len(DB.entries) > 0 {
		t.Log("Database should start out with no entries")
		t.Fail()
	}

	err := AddFile("images/ascendingGradient.png")
	if err != nil {
		t.Logf("Expected no error while adding file.  Got: %v", err)
		t.Fail()
	}

	if len(DB.entries) != 1 {
		t.Log("Database should contain exactly one entry")
		t.Fail()
	}

}

func TestSearchByFile(t *testing.T) {
	AddFile("images/ascendingGradient.png")
	AddFile("images/descendingGradient.png")
	AddFile("images/alternatingGradient.png")
	AddFile("images/gopher1.png")
	AddFile("images/gopher2.png")

	filenames, err := SearchByFile("images/gopher3.png", 5)
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
