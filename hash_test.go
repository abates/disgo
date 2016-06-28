package disgo

import (
	"os"
	"testing"
)

func testHash(filename string, expected PHash, t *testing.T) {
	file, _ := os.Open(filename)
	h, _ := HashFile(file)
	if h != expected {
		t.Logf("Failed to calculate correct hash for image.  Hash was %x expecting %x", h, expected)
		t.Fail()
	}
}

func TestHash(t *testing.T) {
	testHash("images/ascendingGradient.png", 0, t)
	testHash("images/descendingGradient.png", 0xffffffffffffffff, t)
	testHash("images/alternatingGradient.png", 0x00ff00ff00ff00ff, t)
}

func testDistance(v1, v2 PHash, expected int, t *testing.T) {
	d := v1.Distance(v2)
	if d != expected {
		t.Logf("Failed to compute hamming distance.  Expected %d got %d", expected, d)
		t.Fail()
	}

	d = v2.Distance(v1)
	if d != expected {
		t.Logf("Failed to compute hamming distance.  Expected %d got %d", expected, d)
		t.Fail()
	}
}

func TestDistance(t *testing.T) {
	testDistance(0x00, 0xff, 8, t)
	testDistance(0xf0, 0x0f, 8, t)
	testDistance(0xff, 0xfe, 1, t)
}
