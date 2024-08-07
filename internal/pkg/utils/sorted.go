package utils

type TypeSorted int

const (
	Created   TypeSorted = 0
	Updated   TypeSorted = 1
	Visit     TypeSorted = 2
	ABC       TypeSorted = 3
	WordCount TypeSorted = 4
)

type TypeOrder int

const (
	ASC  TypeOrder = 0
	DESC TypeOrder = 1
)

func (t TypeOrder) String() string {
	switch t {
	case DESC:
		return "DESC"
	default:
		return "ASC"
	}
}
