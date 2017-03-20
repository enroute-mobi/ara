package model

type Attributes map[string]string

func NewAttributes() Attributes {
	return make(Attributes)
}

func (attributes Attributes) Set(key string, value string) {
	if value == "" {
		return
	}
	attributes[key] = value
}
