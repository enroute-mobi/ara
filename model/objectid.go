package model

type ObjectIDs map[string]ObjectID

type ObjectID struct {
	kind  string
	value string
}

func NewObjectID(kind, value string) ObjectID {
	return ObjectID{
		kind,
		value,
	}
}

func (objectid ObjectID) Kind() string {
	return objectid.kind
}

func (objectid ObjectID) Value() string {
	return objectid.value
}
