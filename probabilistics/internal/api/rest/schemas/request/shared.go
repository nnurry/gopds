package request_schema

type MetaBody struct {
	Key string `json:"key"`
}

type PairMetaBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
