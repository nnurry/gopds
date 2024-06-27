package cardinal_filter

import (
	"gopds/probabilistics/pkg/models/cardinal"
	"gopds/probabilistics/pkg/models/filter"
	"gopds/probabilistics/pkg/models/meta"
)

type CardinalFilter struct {
	filter   filter.Filter
	cardinal cardinal.Cardinal
	meta     *meta.CardinalFilterMeta
}

func (cf *CardinalFilter) Filter() filter.Filter {
	return cf.filter
}

func (cf *CardinalFilter) Cardinal() cardinal.Cardinal {
	return cf.cardinal
}

func NewCardinalFilter(
	filter filter.Filter, cardinal cardinal.Cardinal, meta *meta.CardinalFilterMeta,
) *CardinalFilter {
	return &CardinalFilter{
		filter:   filter,
		cardinal: cardinal,
		meta:     meta,
	}
}
