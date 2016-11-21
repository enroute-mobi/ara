package model

type ObjectIDs map[string]ObjectID

type ObjectID struct {
	kind  string
	value string
}

func NewObjectID(kind, value string) *ObjectID {
	return &ObjectID{
		kind,
		value,
	}
}

func (objectID *ObjectID) Kind() string {
	return objectID.kind
}

func (objectID *ObjectID) Value() string {
	return objectID.value
}
