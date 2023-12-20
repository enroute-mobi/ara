package model

import (
	"encoding/json"
)

type Reference struct {
	Code *Code  `json:",omitempty"`
	Type string `json:",omitempty"`
}

func NewReference(code Code) *Reference {
	return &Reference{Code: &code}
}

func (reference *Reference) GetSha1() string {
	return reference.Code.HashValue()
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
