package mappers

type Mapper struct {
	methods IMapper
}

type IMapper interface {
	Insert() int
	Delete() int
	Update() int
	SelectOne()
}
