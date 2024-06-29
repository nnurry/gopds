package abstractmeta

type CardinalMeta interface {
	Id() uint
	SetId(uint)
	Key() string
	CardinalType() string
}
