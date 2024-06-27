package cardinal

import "gopds/probabilistics/pkg/models/meta"

type Cardinal interface {
	Meta() *meta.CardinalMeta
	Add([]byte) error
	AddString(string) error
	Cardinality() uint64
}
