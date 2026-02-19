package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	valueReplace = "%{value}"
)

type createCodeUpdaterAttributes struct {
	SourceCodeSpace string `json:"source_code_space"`
	TargetCodeSpace string `json:"target_code_space"`
	TargetPattern   string `json:"target_pattern"`
}

func NewCreateCodeUpdater(m Model, sm *SelectMacro) (updater, error) {
	if !sm.Attributes.Valid {
		return nil, errors.New("empty Attributes")
	}

	var attrs createCodeUpdaterAttributes
	err := json.Unmarshal([]byte(sm.Attributes.String), &attrs)
	if err != nil {
		return nil, fmt.Errorf("can't parse Attributes: %v", err)
	}

	manager := manager(m, sm.ModelType.String)
	if manager == nil {
		return nil, fmt.Errorf("unsupported type %v", sm.Type)
	}

	return func(mi ModelInstance) error {
		_, ok := mi.Code(attrs.TargetCodeSpace)
		if ok {
			return nil
		}

		sourceCode, ok := mi.Code(attrs.SourceCodeSpace)
		if !ok {
			return fmt.Errorf("cannot find source code")
		}

		code := NewCode(attrs.TargetCodeSpace, strings.ReplaceAll(attrs.TargetPattern, valueReplace, sourceCode.Value()))
		if manager.CodeExists(code) {
			return fmt.Errorf("macro CreateCode should create code \"%s\":\"%s\" but it already exists", code.codeSpace, code.value)
		}

		mi.SetCode(code)

		return nil
	}, nil
}

func manager(m Model, t string) CodeValidator {
	switch t {
	case "StopArea":
		return m.StopAreas()
	case "Line":
		return m.Lines()
	case "VehicleJourney":
		return m.VehicleJourneys()
	case "StopVisit":
		return m.StopVisits()
	case "Vehicle":
		return m.Vehicles()
	}
	return nil
}
