package meta

type CardinalMeta struct {
	cardinalType string
}

func (cm *CardinalMeta) CardinalType() string {
	return cm.cardinalType
}

func NewCardinalMeta(cardinalType string) *CardinalMeta {
	return &CardinalMeta{
		cardinalType: cardinalType,
	}
}
