package model

type Startable interface {
	Start()
}

type Stopable interface {
	Stop()
}
