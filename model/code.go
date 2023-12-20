package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
)

type Codes map[string]Code

func NewCodesFromMap(codeMap map[string]string) (codes Codes) {
	codes = make(Codes)
	for key, value := range codeMap {
		codes[key] = NewCode(key, value)
	}
	return codes
}

func (identifiers Codes) Empty() bool {
	return len(identifiers) == 0
}

func (identifiers Codes) UnmarshalJSON(text []byte) error {
	var definitions map[string]string
	if err := json.Unmarshal(text, &definitions); err != nil {
		return err
	}
	for key, value := range definitions {
		identifiers[key] = NewCode(key, value)
	}
	return nil
}

func (identifiers Codes) MarshalJSON() ([]byte, error) {
	aux := map[string]string{}

	for codeSpace, code := range identifiers {
		aux[codeSpace] = code.Value()
	}

	return json.Marshal(aux)
}

func (identifiers Codes) ToSlice() (objs []string) {
	for _, obj := range identifiers {
		objs = append(objs, obj.String())
	}
	return
}

type Code struct {
	codeSpace string
	value     string
}

func NewCode(codeSpace, value string) Code {
	return Code{
		codeSpace,
		value,
	}
}

func (code Code) CodeSpace() string {
	return code.codeSpace
}

func (code Code) Value() string {
	return code.value
}

func (code Code) HashValue() string {
	hasher := sha1.New() // oui, on sait
	hasher.Write([]byte(code.Value()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (code Code) String() string {
	return fmt.Sprintf("%s:%s", code.codeSpace, code.value)
}

func (code *Code) SetValue(toset string) {
	code.value = toset
}

func (code *Code) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		code.codeSpace: code.value,
	})
}

func (code *Code) UnmarshalJSON(data []byte) error {
	var aux map[string]string

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux == nil {
		return nil
	}

	if len(aux) > 1 {
		return errors.New("Code should look like CODESPACE:VALUE")
	}

	for codeSpace, value := range aux {
		code.codeSpace = codeSpace
		code.value = value
	}

	return nil
}

type CodeConsumerInterface interface {
	Code(string) (Code, bool)
	CodeWithFallback([]string) (Code, bool)
	Codes() Codes
	CodesResponse() map[string]string
	SetCode(Code)
	CodeSlice() []string
}

type CodeConsumer struct {
	codes Codes
}

func (consumer *CodeConsumer) Copy() CodeConsumer {
	o := CodeConsumer{
		codes: make(Codes),
	}
	for k, v := range consumer.codes {
		o.codes[k] = v
	}
	return o
}

func (consumer *CodeConsumer) Code(codeSpace string) (Code, bool) {
	code, ok := consumer.codes[codeSpace]
	if ok {
		return code, true
	}
	return Code{}, false
}

func (consumer *CodeConsumer) CodeWithFallback(codeSpaces []string) (Code, bool) {
	for i := range codeSpaces {
		code, ok := consumer.codes[codeSpaces[i]]
		if ok {
			return code, true
		}
	}
	return Code{}, false
}

func (consumer *CodeConsumer) SetCode(code Code) {
	consumer.codes[code.CodeSpace()] = code
}

func (consumer *CodeConsumer) Codes() Codes {
	return consumer.codes
}

func (consumer *CodeConsumer) CodesResponse() map[string]string {
	codes := make(map[string]string)
	for _, code := range consumer.codes {
		codes[code.CodeSpace()] = code.Value()
	}
	return codes
}

func (consumer *CodeConsumer) CodeSlice() (objs []string) {
	return consumer.codes.ToSlice()
}
