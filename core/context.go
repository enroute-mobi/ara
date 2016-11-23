package core

import "github.com/af83/edwig/model"

type Context map[string]interface{}

func (context *Context) IsDefined(key string) bool {
	_, ok := (*context)[key]
	return ok
}

func (context *Context) Value(key string) interface{} {
	return (*context)[key]
}

func (context *Context) SetValue(key string, value interface{}) {
	(*context)[key] = value
}

func (context *Context) Close() {
	for _, contextElement := range *context {
		_, ok := contextElement.(model.Stopable)
		if ok {
			contextElement.(model.Stopable).Stop()
		}
	}
}
