package model

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
