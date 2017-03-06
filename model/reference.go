package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
)

type Reference struct {
	ObjectId *ObjectID
	Id       string
}

func (reference *Reference) ChecksumObjId() {
	hasher := sha1.New() // oui, on sait
	hasher.Write([]byte(reference.ObjectId.Value()))
	sha := fmt.Sprintf("%x", hasher.Sum(nil))
	reference.ObjectId.SetValue(sha)
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
