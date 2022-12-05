package sxml

type Bool struct {
	Value   bool
	Defined bool // Valid is true if Bool is not NULL
}

func (b *Bool) SetValue(value bool) {
	b.Value = value
	b.Defined = true
}

type Int struct {
	Value   int
	Defined bool
}

func (i *Int) SetValue(value int) {
	i.Value = value
	i.Defined = true
}

func (i *Int) SetValueWithDefault(value, d int) {
	i.Defined = true
	if value == 0 {
		i.Value = d
		return
	}
	i.Value = value
}
