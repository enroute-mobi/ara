package model

import "encoding/json"

type ObjectIDs map[string]ObjectID

func NewObjectIDsFromMap(objectIdMap map[string]string) (objectids ObjectIDs) {
	objectids = make(ObjectIDs)
	for key, value := range objectIdMap {
		objectids[key] = NewObjectID(key, value)
	}
	return objectids
}

/*func (identifiers ObjectIDs) UnmarshalJSON(text []byte) error {
	var definitions map[string]string
	if err := json.Unmarshal(text, &definitions); err != nil {
		return err
	}
	for key, value := range definitions {
		identifiers[key] = NewObjectID(key, value)
	}
	return nil
} */

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

func (objectid *ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		objectid.kind: objectid.value,
	})
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

func (consumer *ObjectIDConsumer) ObjectIDsResponse() map[string]string {
	objectIds := make(map[string]string)
	for _, object := range consumer.objectids {
		objectIds[object.Kind()] = object.Value()
	}
	return objectIds
}
