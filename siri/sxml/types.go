package sxml

type Bool struct {
	Value   bool
	Defined bool // Valid is true if Bool is not NULL
}

func (b *Bool) SetValue(value bool) {
	b.Value = value
	b.Defined = true
}

func (b *Bool) Parse(value string) {
	b.SetValue(value == "true" || value == "TRUE" || value == "1")
}
