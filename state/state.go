package state

type Startable interface {
	Start()
}

type Stopable interface {
	Stop()
}
