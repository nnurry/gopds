package probabilistic

import (
	"gopds/probabilistics/pkg/models/cardinal"
	"gopds/probabilistics/pkg/models/filter"
	"gopds/probabilistics/pkg/models/meta"
)

type Probabilistic struct {
	filter   filter.Filter
	cardinal cardinal.Cardinal
	meta     *meta.ProbabilisticMeta
}

func (p *Probabilistic) Filter() filter.Filter {
	return p.filter
}

func (p *Probabilistic) Cardinal() cardinal.Cardinal {
	return p.cardinal
}

func (p *Probabilistic) Meta() *meta.ProbabilisticMeta {
	return p.meta
}

func NewProbabilistic(
	filter filter.Filter, cardinal cardinal.Cardinal, meta *meta.ProbabilisticMeta) *Probabilistic {
	return &Probabilistic{
		filter:   filter,
		cardinal: cardinal,
		meta:     meta,
	}
}
