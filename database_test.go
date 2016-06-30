package disgo

import (
	"math/rand"
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
	db := New()
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
	for _, i := range []Index{NewRadixIndex(), NewLinearIndex()} {
		db := NewDB(i)
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
}

var benchmarkHashes []PHash

func getHashes(numHashes int) []PHash {
	if benchmarkHashes == nil {
		benchmarkHashes = make([]PHash, 0)
	}

	for len(benchmarkHashes) < numHashes {
		randomNumber := PHash(rand.Int63())
		if rand.NormFloat64() >= 0 {
			randomNumber = randomNumber | 0x8000000000000000
		}
		benchmarkHashes = append(benchmarkHashes, randomNumber)
	}
	return benchmarkHashes
}

func benchmarkAdd(b *testing.B, index Index, numToAdd int) {
	hashes := getHashes(numToAdd)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		db := NewDB(index)
		for i := 0; i < numToAdd; i++ {
			db.Add("filename", hashes[i])
		}
	}
}

func benchmarkSearch(b *testing.B, index Index, numToSearch int) {
	hashes := getHashes(numToSearch)

	db := NewDB(index)
	for _, hash := range hashes {
		db.Add("filename", hash)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < numToSearch; i++ {
			db.Search(hashes[i], 5)
		}
	}
}

func BenchmarkLinearIndexAdd10(b *testing.B)    { benchmarkAdd(b, NewLinearIndex(), 10) }
func BenchmarkLinearIndexAdd100(b *testing.B)   { benchmarkAdd(b, NewLinearIndex(), 100) }
func BenchmarkLinearIndexAdd1000(b *testing.B)  { benchmarkAdd(b, NewLinearIndex(), 1000) }
func BenchmarkLinearIndexAdd10000(b *testing.B) { benchmarkAdd(b, NewLinearIndex(), 10000) }

func BenchmarkLinearIndexSearch10(b *testing.B)    { benchmarkSearch(b, NewLinearIndex(), 10) }
func BenchmarkLinearIndexSearch100(b *testing.B)   { benchmarkSearch(b, NewLinearIndex(), 100) }
func BenchmarkLinearIndexSearch1000(b *testing.B)  { benchmarkSearch(b, NewLinearIndex(), 1000) }
func BenchmarkLinearIndexSearch10000(b *testing.B) { benchmarkSearch(b, NewLinearIndex(), 10000) }

func BenchmarkRadixIndexAdd10(b *testing.B)    { benchmarkAdd(b, NewRadixIndex(), 10) }
func BenchmarkRadixIndexAdd100(b *testing.B)   { benchmarkAdd(b, NewRadixIndex(), 100) }
func BenchmarkRadixIndexAdd1000(b *testing.B)  { benchmarkAdd(b, NewRadixIndex(), 1000) }
func BenchmarkRadixIndexAdd10000(b *testing.B) { benchmarkAdd(b, NewRadixIndex(), 10000) }

func BenchmarkRadixIndexSearch10(b *testing.B)     { benchmarkSearch(b, NewRadixIndex(), 10) }
func BenchmarkRadixIndexSearch100(b *testing.B)    { benchmarkSearch(b, NewRadixIndex(), 100) }
func BenchmarkRadixIndexSearch1000(b *testing.B)   { benchmarkSearch(b, NewRadixIndex(), 1000) }
func BenchmarkRadixIndexSearch10000(b *testing.B)  { benchmarkSearch(b, NewRadixIndex(), 10000) }
func BenchmarkRadixIndexSearch100000(b *testing.B) { benchmarkSearch(b, NewRadixIndex(), 100000) }
