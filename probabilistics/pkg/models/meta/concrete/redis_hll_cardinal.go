package concretemeta

type RedisHyperLogLogMeta struct {
	id           uint
	cardinalType string
	key          string
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

func (m *RedisHyperLogLogMeta) Key() string {
	return m.key
}

func NewRedisHLLMeta(cardinalType string, key string) *RedisHyperLogLogMeta {
	return &RedisHyperLogLogMeta{
		cardinalType: cardinalType,
		key:          key,
	}
}
