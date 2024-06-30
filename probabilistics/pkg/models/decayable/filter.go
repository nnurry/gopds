package decayable

import (
	abstractfilter "github.com/nnurry/gopds/probabilistics/pkg/models/filter/abstract"
	concretemeta "github.com/nnurry/gopds/probabilistics/pkg/models/meta/concrete"
)

type Filter struct {
	core abstractfilter.Filter
	meta *concretemeta.DecayableMeta
}

func (p *Filter) Core() abstractfilter.Filter {
	return p.core
}

func (p *Filter) SetCore(core abstractfilter.Filter) {
	p.core = core
}

func (p *Filter) Meta() *concretemeta.DecayableMeta {
	return p.meta
}

func (p *Filter) SetMeta(m *concretemeta.DecayableMeta) {
	p.meta = m
}

func (p *Filter) Add(value []byte) {
	if err := p.core.Add(value); err != nil {
		panic(err)
	}
}

func (p *Filter) AddString(value string) {
	if err := p.core.AddString(value); err != nil {
		panic(err)
	}
}
