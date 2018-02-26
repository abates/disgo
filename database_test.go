package disgo

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"reflect"
	"testing"
)

type testIndex struct {
	err     error
	matches []PHash
}

func (ti *testIndex) Insert(PHash) error                 { return ti.err }
func (ti *testIndex) Search(PHash, int) ([]PHash, error) { return ti.matches, ti.err }

func newTestIndex() *testIndex {
	return &testIndex{}
}

func TestDBAdd(t *testing.T) {
	tests := []struct {
		img         image.Image
		expectedErr error
	}{
		{image.NewRGBA(image.Rect(0, 0, 10, 10)), nil},
		{image.NewRGBA(image.Rect(0, 0, 10, 10)), fmt.Errorf("Just some test error")},
	}

	for i, test := range tests {
		db := NewDB(newTestIndex())
		db.hasher = func(image.Image) (PHash, error) { return PHash(0), test.expectedErr }

		_, err := db.Add(test.img)
		if err != test.expectedErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedErr, err)
		}
	}
}

func TestDBAddFile(t *testing.T) {
	tests := []struct {
		expectedErr error
	}{
		{nil},
		{fmt.Errorf("Just some test error")},
	}

	for i, test := range tests {
		db := NewDB(newTestIndex())
		db.hasher = func(image.Image) (PHash, error) { return PHash(0), test.expectedErr }
		image := image.NewAlpha(image.Rect(0, 0, 1, 1))
		buf := bytes.NewBuffer([]byte{})
		png.Encode(buf, image)
		_, err := db.AddFile(buf)
		if err != test.expectedErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedErr, err)
		}
	}
}

func TestDBSearch(t *testing.T) {
	tests := []struct {
		img               image.Image
		expectedMatches   []PHash
		expectedHashErr   error
		expectedSearchErr error
	}{
		{image.NewRGBA(image.Rect(0, 0, 10, 10)), []PHash{}, nil, nil},
		{image.NewRGBA(image.Rect(0, 0, 10, 10)), []PHash{}, nil, fmt.Errorf("Some test error")},
	}

	for i, test := range tests {
		testIndex := newTestIndex()
		testIndex.err = test.expectedSearchErr
		testIndex.matches = test.expectedMatches
		db := NewDB(testIndex)
		db.hasher = func(image.Image) (PHash, error) { return PHash(0), test.expectedHashErr }

		matches, err := db.Search(test.img, 0)
		if test.expectedHashErr != nil && err != test.expectedHashErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedHashErr, err)
		} else if test.expectedSearchErr != nil && err != test.expectedSearchErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedSearchErr, err)
		}

		if !reflect.DeepEqual(test.expectedMatches, matches) {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedMatches, matches)
		}
	}
}

func TestDBSearchByFile(t *testing.T) {
	tests := []struct {
		expectedMatches   []PHash
		expectedHashErr   error
		expectedSearchErr error
	}{
		{[]PHash{}, nil, nil},
		{[]PHash{}, nil, fmt.Errorf("Some test error")},
	}

	for i, test := range tests {
		testIndex := newTestIndex()
		testIndex.err = test.expectedSearchErr
		testIndex.matches = test.expectedMatches
		db := NewDB(testIndex)
		db.hasher = func(image.Image) (PHash, error) { return PHash(0), test.expectedHashErr }

		image := image.NewAlpha(image.Rect(0, 0, 1, 1))
		buf := bytes.NewBuffer([]byte{})
		png.Encode(buf, image)
		matches, err := db.SearchByFile(buf, 0)
		if test.expectedHashErr != nil && err != test.expectedHashErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedHashErr, err)
		} else if test.expectedSearchErr != nil && err != test.expectedSearchErr {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedSearchErr, err)
		}

		if !reflect.DeepEqual(test.expectedMatches, matches) {
			t.Errorf("tests[%d] expected %v got %v", i, test.expectedMatches, matches)
		}
	}
}

/*
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

func BenchmarkRadixIndexSearch10(b *testing.B)    { benchmarkSearch(b, NewRadixIndex(), 10) }
func BenchmarkRadixIndexSearch100(b *testing.B)   { benchmarkSearch(b, NewRadixIndex(), 100) }
func BenchmarkRadixIndexSearch1000(b *testing.B)  { benchmarkSearch(b, NewRadixIndex(), 1000) }
func BenchmarkRadixIndexSearch10000(b *testing.B) { benchmarkSearch(b, NewRadixIndex(), 10000) }
*/
