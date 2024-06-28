package abstractcardinal

import abstractmeta "gopds/probabilistics/pkg/models/meta/abstract"

type Cardinal interface {
	Meta() abstractmeta.CardinalMeta
	Add([]byte) error
	AddString(string) error
	Cardinality() uint64
	Deserialize([]byte) error
	Serialize() []byte
}
