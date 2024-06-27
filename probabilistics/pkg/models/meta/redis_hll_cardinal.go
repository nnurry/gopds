package meta

type RedisHyperLogLogMeta struct {
	CardinalMeta
	pfKey string
}

func (cm *RedisHyperLogLogMeta) PFKey() string {
	return cm.pfKey
}

func NewRedisHLLMeta(cardinalType string, pfKey string) *RedisHyperLogLogMeta {
	return &RedisHyperLogLogMeta{
		CardinalMeta: CardinalMeta{cardinalType: cardinalType},
		pfKey:        pfKey,
	}
}
