package decayable

import (
	abstractcardinal "gopds/probabilistics/pkg/models/cardinal/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"
)

type Cardinal struct {
	core abstractcardinal.Cardinal
	meta *concretemeta.DecayableMeta
}

func (p *Cardinal) Core() abstractcardinal.Cardinal {
	return p.core
}

func (p *Cardinal) SetCore(core abstractcardinal.Cardinal) {
	p.core = core
}

func (p *Cardinal) Meta() *concretemeta.DecayableMeta {
	return p.meta
}

func (p *Cardinal) SetMeta(m *concretemeta.DecayableMeta) {
	p.meta = m
}

func (p *Cardinal) Add(value []byte) {
	if err := p.core.Add(value); err != nil {
		panic(err)
	}
}

func (p *Cardinal) AddString(value string) {
	if err := p.core.AddString(value); err != nil {
		panic(err)
	}
}
