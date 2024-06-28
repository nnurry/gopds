package concretemeta

type StandardHyperLogLogMeta struct {
	id           uint
	cardinalType string
	pfKey        string
}

func (m *StandardHyperLogLogMeta) Id() uint {
	return m.id
}

func (m *StandardHyperLogLogMeta) SetId(id uint) {
	m.id = id
}

func (m *StandardHyperLogLogMeta) CardinalType() string {
	return m.cardinalType
}

func (m *StandardHyperLogLogMeta) PFKey() string {
	return m.pfKey
}

func NewStandardHLLMeta(cardinalType string, pfKey string) *StandardHyperLogLogMeta {
	return &StandardHyperLogLogMeta{
		cardinalType: cardinalType,
		pfKey:        pfKey,
	}
}
