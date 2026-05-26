package model

import (
	"encoding/json"
	"time"
)

type LineId string

type Line struct {
	Collectable
	model      Model
	References References
	CodeConsumer
	RawAttributes     RawAttributes
	id                LineId
	ReferentId        LineId `json:",omitempty"`
	Name              string `json:",omitempty"`
	Number            string `json:",omitempty"`
	Origin            string `json:",omitempty"`
	CollectSituations bool
}

func NewLine(model Model) *Line {
	line := &Line{
		model:         model,
		RawAttributes: NewRawAttributes(),
		References:    NewReferences(),
	}

	line.InitCodes()
	return line
}

func (line *Line) ModelId() string {
	return string(line.id)
}

func (line *Line) GetName() string {
	return line.Name
}

func (line *Line) copy() *Line {
	l := *line
	l.RawAttributes = line.RawAttributes.Copy()
	l.References = line.References.Copy()
	return &l
}

func (line *Line) Id() LineId {
	return line.id
}

func (line *Line) SetId(id LineId) {
	line.id = id
}

func (line *Line) SetOrigin(origin string) {
	line.Origin = origin
}

func (line *Line) MarshalJSON() ([]byte, error) {
	type Alias Line
	aux := struct {
		*Alias
		Codes         Codes                `json:",omitempty"`
		NextCollectAt *time.Time           `json:",omitempty"`
		CollectedAt   *time.Time           `json:",omitempty"`
		RawAttributes RawAttributes        `json:",omitempty"`
		References    map[string]Reference `json:",omitempty"`
		Id            LineId
	}{
		Id:    line.id,
		Alias: (*Alias)(line),
	}

	if !line.Codes().Empty() {
		aux.Codes = line.Codes()
	}
	if !line.nextCollectAt.IsZero() {
		aux.NextCollectAt = &line.nextCollectAt
	}
	if !line.collectedAt.IsZero() {
		aux.CollectedAt = &line.collectedAt
	}
	if !line.RawAttributes.IsEmpty() {
		aux.RawAttributes = line.RawAttributes
	}

	if !line.References.IsEmpty() {
		aux.References = line.References.GetReferences()
	}

	return json.Marshal(&aux)
}

func (line *Line) UnmarshalJSON(data []byte) error {
	type Alias Line

	aux := &struct {
		Codes      map[string]string
		References map[string]Reference
		Id         string
		*Alias
	}{
		Alias: (*Alias)(line),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		line.SetCodesFromMap(aux.Codes)
	}

	if aux.References != nil {
		line.References.SetReferences(aux.References)
	}

	if aux.Id != "" && line.id == "" {
		line.id = LineId(aux.Id)
	}
	return nil
}

func (line *Line) Referent() (*Line, bool) {
	if line.ReferentId == "" {
		return nil, false
	}
	return line.model.Lines().Find(line.ReferentId)
}

func (line *Line) ReferentOrSelfCode(codeSpace string) (Code, bool) {
	if line.ReferentId != "" {
		if refCode, ok := line.model.Lines().FindCode(line.ReferentId, codeSpace); ok {
			return refCode, true
		}
	}

	code, ok := line.Code(codeSpace)
	if ok {
		return code, true
	}
	return Code{}, false
}

/*
Returns true if we need to send the Line in a Discovery

We only send the Line if it has no referent with a correct codeSpace.
If that's the case, we'll send the Referent instead
*/
func (line *Line) DiscoveryCode(codeSpace string) (Code, bool) {
	if line.ReferentId != "" {
		if _, ok := line.model.Lines().FindCode(line.ReferentId, codeSpace); ok {
			return Code{}, false
		}
	}

	code, ok := line.Code(codeSpace)
	if ok {
		return code, true
	}
	return Code{}, false
}

func (line *Line) Attribute(key string) (string, bool) {
	value, present := line.RawAttributes[key]
	return value, present
}

func (line *Line) Reference(key string) (Reference, bool) {
	value, present := line.References.Get(key)
	return value, present
}

func (line *Line) Save() bool {
	return line.model.Lines().Save(line)
}
