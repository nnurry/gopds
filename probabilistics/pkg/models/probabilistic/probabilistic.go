package probabilistic

import (
	abstractcardinal "gopds/probabilistics/pkg/models/cardinal/abstract"
	abstractfilter "gopds/probabilistics/pkg/models/filter/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"
)

type Probabilistic struct {
	filter   abstractfilter.Filter
	cardinal abstractcardinal.Cardinal
	meta     *concretemeta.ProbabilisticMeta
}

func (p *Probabilistic) Filter() abstractfilter.Filter {
	return p.filter
}

func (p *Probabilistic) Cardinal() abstractcardinal.Cardinal {
	return p.cardinal
}

func (p *Probabilistic) Meta() *concretemeta.ProbabilisticMeta {
	return p.meta
}

func (p *Probabilistic) SetFilter(f abstractfilter.Filter) {
	p.filter = f
}

func (p *Probabilistic) SetCardinal(c abstractcardinal.Cardinal) {
	p.cardinal = c
}

func (p *Probabilistic) SetMeta(m *concretemeta.ProbabilisticMeta) {
	p.meta = m
}

func (p *Probabilistic) Add(value []byte) {
	var err error
	err = p.Filter().Add(value)
	if err != nil {
		panic(err)
	}
	err = p.Cardinal().Add(value)
	if err != nil {
		panic(err)
	}
}

func (p *Probabilistic) AddString(value string) {
	var err error
	err = p.Filter().AddString(value)
	if err != nil {
		panic(err)
	}
	err = p.Cardinal().AddString(value)
	if err != nil {
		panic(err)
	}
}

func NewProbabilistic(
	filter abstractfilter.Filter,
	cardinal abstractcardinal.Cardinal,
	meta *concretemeta.ProbabilisticMeta) *Probabilistic {
	return &Probabilistic{
		filter:   filter,
		cardinal: cardinal,
		meta:     meta,
	}
}
