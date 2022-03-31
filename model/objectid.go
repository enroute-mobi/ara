package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
)

type ObjectIDs map[string]ObjectID

func NewObjectIDsFromMap(objectIdMap map[string]string) (objectids ObjectIDs) {
	objectids = make(ObjectIDs)
	for key, value := range objectIdMap {
		objectids[key] = NewObjectID(key, value)
	}
	return objectids
}

func (identifiers ObjectIDs) Empty() bool {
	return len(identifiers) == 0
}

func (identifiers ObjectIDs) UnmarshalJSON(text []byte) error {
	var definitions map[string]string
	if err := json.Unmarshal(text, &definitions); err != nil {
		return err
	}
	for key, value := range definitions {
		identifiers[key] = NewObjectID(key, value)
	}
	return nil
}

func (identifiers ObjectIDs) MarshalJSON() ([]byte, error) {
	aux := map[string]string{}

	for kind, objectid := range identifiers {
		aux[kind] = objectid.Value()
	}

	return json.Marshal(aux)
}

func (identifiers ObjectIDs) ToSlice() (objs []string) {
	for _, obj := range identifiers {
		objs = append(objs, obj.String())
	}
	return
}

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

func (objectid ObjectID) HashValue() string {
	hasher := sha1.New() // oui, on sait
	hasher.Write([]byte(objectid.Value()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (objectid ObjectID) String() string {
	return fmt.Sprintf("%s:%s", objectid.kind, objectid.value)
}

func (objectid *ObjectID) SetValue(toset string) {
	objectid.value = toset
}

func (objectid *ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		objectid.kind: objectid.value,
	})
}

func (objectid *ObjectID) UnmarshalJSON(data []byte) error {
	var aux map[string]string

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux == nil {
		return nil
	}

	if len(aux) > 1 {
		return errors.New("ObjectID should look like KIND:VALUE")
	}

	for kind, value := range aux {
		objectid.kind = kind
		objectid.value = value
	}

	return nil
}

type ObjectIDConsumerInterface interface {
	ObjectID(string) (ObjectID, bool)
	ObjectIDWithFallback([]string) (ObjectID, bool)
	ObjectIDs() ObjectIDs
	ObjectIDsResponse() map[string]string
	SetObjectID(ObjectID)
	ObjectIDSlice() []string
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

func (consumer *ObjectIDConsumer) ObjectIDWithFallback(kinds []string) (ObjectID, bool) {
	for i := range kinds {
		objectid, ok := consumer.objectids[kinds[i]]
		if ok {
			return objectid, true
		}
	}
	return ObjectID{}, false
}

func (consumer *ObjectIDConsumer) SetObjectID(objectid ObjectID) {
	consumer.objectids[objectid.Kind()] = objectid
}

func (consumer *ObjectIDConsumer) ObjectIDs() ObjectIDs {
	return consumer.objectids
}

func (consumer *ObjectIDConsumer) ObjectIDsResponse() map[string]string {
	objectIds := make(map[string]string)
	for _, object := range consumer.objectids {
		objectIds[object.Kind()] = object.Value()
	}
	return objectIds
}

func (consumer *ObjectIDConsumer) ObjectIDSlice() (objs []string) {
	return consumer.objectids.ToSlice()
}
