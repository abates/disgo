package disgo

type entry struct {
	filenames []string
}

type db struct {
	entries map[PHash]*entry
}

var DB db

func init() {
	DB.entries = make(map[PHash]*entry)
}

func AddFile(filename string) error {
	h, err := HashFile(filename)
	if err != nil {
		return err
	}
	e, found := DB.entries[h]
	if !found {
		e = new(entry)
		DB.entries[h] = e
	}
	e.filenames = append(e.filenames, filename)
	return nil
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
			results = append(results, e.filenames...)
		}
	}

	return results, nil
}
