package disgo

type LinearIndex struct {
	entries map[PHash]bool
}

func NewLinearIndex() *LinearIndex {
	i := new(LinearIndex)
	i.entries = make(map[PHash]bool)
	return i
}

func (i *LinearIndex) Insert(phash PHash) error {
	i.entries[phash] = true
	return nil
}

func (i *LinearIndex) Search(phash PHash, maxDistance int) ([]PHash, error) {
	var results []PHash

	// look for existing entry within maxDistance of the hash
	for p, _ := range i.entries {
		if p.Distance(phash) <= maxDistance {
			results = append(results, p)
		}
	}

	return results, nil
}
