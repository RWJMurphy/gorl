package gorl

// Item is any carryable game thing
type Item interface {
	Feature
	Weight() int
}

type item struct {
	feature
	weight int
}

func NewItem(name string, char rune, weight int) Item {
	i := &item{
		*NewFeature(name, char).(*feature),
		weight,
	}
	i.flags |= FlagCrossable
	return i
}

func (i *item) Weight() int {
	return i.weight
}
