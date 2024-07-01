package request_schema

type FilterBody struct {
	Type           string  `json:"type"`
	MaxCardinality uint    `json:"max_cardinality"`
	ErrorRate      float64 `json:"error_rate"`
}

type FilterCreateBody struct {
	Meta   MetaBody   `json:"meta"`
	Filter FilterBody `json:"filter"`
}

type FilterExistsBody struct {
	Meta   PairMetaBody `json:"meta"`
	Filter FilterBody   `json:"filter"`
}

type FilterAddBody struct {
	Meta   PairMetaBody `json:"meta"`
	Filter FilterBody   `json:"filter"`
}
