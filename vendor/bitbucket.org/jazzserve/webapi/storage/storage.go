package storage

type SortItem interface {
	Field() string
	SetField(val string)
	Order() string
}
