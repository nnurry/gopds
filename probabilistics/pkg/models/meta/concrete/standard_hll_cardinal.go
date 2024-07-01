package concretemeta

type StandardHyperLogLogMeta struct {
	id           uint
	key          string
	cardinalType string
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

func (m *StandardHyperLogLogMeta) Key() string {
	return m.key
}

func NewStandardHLLMeta(key string) *StandardHyperLogLogMeta {
	return &StandardHyperLogLogMeta{
		cardinalType: "STANDARD_HLL",
		key:          key,
	}
}
