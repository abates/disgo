package disgo

type SearchCriteria struct {
	Hash     PHash `json:"hash"`
	Distance uint  `json:"distance"`
}

type ImageInfo struct {
	Hash     PHash  `json:"hash"`
	Location string `json:"location"`
}
