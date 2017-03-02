package model

import (
	"encoding/json"
	"errors"
)

type Reference struct {
	ObjectId *ObjectID
	Id       string
}

func (reference *Reference) UnmarshalJSON(data []byte) error {

	aux := &struct {
		ObjectId map[string]string
		Id       string
	}{}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if len(aux.ObjectId) != 1 {
		return errors.New("ObjectID should look like KIND:VALUE")
	}

	for kind, _ := range aux.ObjectId {
		ObjectIdCPY := NewObjectID(kind, aux.ObjectId[kind])
		reference.ObjectId = &ObjectIdCPY
	}
	reference.Id = aux.Id
	return nil
}
