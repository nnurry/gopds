package abstractmeta

type CardinalMeta interface {
	Id() uint
	SetId(uint)
	CardinalType() string
}
