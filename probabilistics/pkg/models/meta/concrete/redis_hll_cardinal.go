package concretemeta

type RedisHyperLogLogMeta struct {
	id           uint
	cardinalType string
	pfKey        string
}

func (m *RedisHyperLogLogMeta) Id() uint {
	return m.id
}

func (m *RedisHyperLogLogMeta) SetId(id uint) {
	m.id = id
}

func (m *RedisHyperLogLogMeta) CardinalType() string {
	return m.cardinalType
}

func (m *RedisHyperLogLogMeta) PFKey() string {
	return m.pfKey
}

func NewRedisHLLMeta(cardinalType string, pfKey string) *RedisHyperLogLogMeta {
	return &RedisHyperLogLogMeta{
		cardinalType: cardinalType,
		pfKey:        pfKey,
	}
}
