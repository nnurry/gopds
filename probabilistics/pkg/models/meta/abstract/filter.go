package abstractmeta

type FilterMeta interface {
	Id() uint
	SetId(uint)
	Key() string
	FilterType() string
	MaxCard() uint
	MaxFp() float64
	HashFuncNum() uint
	HashFuncType() string
}
