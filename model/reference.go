package model

import "encoding/json"

type Reference struct {
	ObjectId *ObjectID `json:",omitempty"`
	Id       string    `json:",omitempty"`
	Type     string    `json:",omitempty"`
}

func NewReference(objectId ObjectID) *Reference {
	return &Reference{ObjectId: &objectId}
}

func (reference *Reference) GetSha1() string {
	return reference.ObjectId.HashValue()
}

func (reference *Reference) UnmarshalJSON(data []byte) error {
	type Alias Reference
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(reference),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	return nil
}
