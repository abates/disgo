package disgo

type LinearIndex struct {
	entries map[PHash]bool
}

func NewLinearIndex() *LinearIndex {
	li := new(LinearIndex)
	li.entries = make(map[PHash]bool)
	return li
}

func (li *LinearIndex) Insert(phash PHash) error {
	li.entries[phash] = true
	return nil
}

func (li *LinearIndex) Search(phash PHash, maxDistance int) ([]PHash, error) {
	var results []PHash

	// look for existing entry within maxDistance of the hash
	for p := range li.entries {
		if p.Distance(phash) <= maxDistance {
			results = append(results, p)
		}
	}

	return results, nil
}

func (li *LinearIndex) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (li *LinearIndex) UnmarshalBinary(buf []byte) error {
	return nil
}
