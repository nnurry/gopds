package meta

type FilterMeta struct {
	filterType   string
	maxCard      uint
	maxFp        float64
	hashFuncNum  uint
	hashFuncType string
}

func (fm *FilterMeta) FilterType() string {
	return fm.filterType
}
func (fm *FilterMeta) MaxCard() uint {
	return fm.maxCard
}
func (fm *FilterMeta) MaxFp() float64 {
	return fm.maxFp
}
func (fm *FilterMeta) HashFuncNum() uint {
	return fm.hashFuncNum
}
func (fm *FilterMeta) HashFuncType() string {
	return fm.hashFuncType
}

func NewFilterMeta(
	filterType string, maxCard uint, maxFp float64,
	hashFuncNum uint, hashFuncType string,
) *FilterMeta {
	return &FilterMeta{
		filterType:   filterType,
		maxCard:      maxCard,
		maxFp:        maxFp,
		hashFuncNum:  hashFuncNum,
		hashFuncType: hashFuncType,
	}
}
