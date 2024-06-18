package model

import "golang.org/x/exp/maps"

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

func (attributes Attributes) IsEmpty() bool {
	return len(attributes) == 0
}

func (attributes Attributes) Copy() (c Attributes) {
	c = maps.Clone(attributes)

	return
}
