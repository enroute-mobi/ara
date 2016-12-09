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

type ObjectIDConsumer struct {
	objectids ObjectIDs
}

func (consumer *ObjectIDConsumer) ObjectID(kind string) (ObjectID, bool) {
	objectid, ok := consumer.objectids[kind]
	if ok {
		return objectid, true
	}
	return ObjectID{}, false
}

func (consumer *ObjectIDConsumer) SetObjectID(objectid ObjectID) {
	consumer.objectids[objectid.Kind()] = objectid
}

func (consumer *ObjectIDConsumer) ObjectIDs() (objectidArray []ObjectID) {
	for _, objectid := range consumer.objectids {
		objectidArray = append(objectidArray, objectid)
	}
	return
}
