package filter

import "gopds/probabilistics/pkg/models/meta"

type Filter interface {
	Meta() *meta.FilterMeta
	Add([]byte) error
	AddString(string) error
	Exists([]byte) bool
}
