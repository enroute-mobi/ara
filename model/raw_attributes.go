package model

import "golang.org/x/exp/maps"

type RawAttributes map[string]string

func NewRawAttributes() RawAttributes {
	return make(RawAttributes)
}

func (attributes RawAttributes) Set(key string, value string) {
	if value == "" {
		return
	}
	attributes[key] = value
}

func (attributes RawAttributes) IsEmpty() bool {
	return len(attributes) == 0
}

func (attributes RawAttributes) Copy() (c RawAttributes) {
	c = maps.Clone(attributes)

	return
}
