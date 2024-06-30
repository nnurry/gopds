package abstractfilter

import abstractmeta "github.com/nnurry/gopds/probabilistics/pkg/models/meta/abstract"

type Filter interface {
	Meta() abstractmeta.FilterMeta
	Add([]byte) error
	AddString(string) error
	Exists([]byte) bool
	Deserialize([]byte) error
	Serialize() []byte
}
