package request_schema

type CardinalBody struct {
	Type string `json:"type"`
}

type CardinalCreateBody struct {
	Meta     MetaBody     `json:"meta"`
	Cardinal CardinalBody `json:"cardinal"`
}

type CardinalCardBody struct {
	Meta     MetaBody     `json:"meta"`
	Cardinal CardinalBody `json:"cardinal"`
}

type CardinalAddBody struct {
	Meta     PairMetaBody `json:"meta"`
	Cardinal CardinalBody `json:"cardinal"`
}
