package core

import "strings"

var defaultIdentifierGenerators = map[string]string{
	"message_identifier":             "%{uuid}",
	"response_message_identifier":    "%{uuid}",
	"data_frame_identifier":          "%{id}",
	"reference_identifier":           "%{type}:%{default}",
	"reference_stop_area_identifier": "%{default}",
}

type IdentifierGenerator struct {
	formatString string
}

type IdentifierAttributes struct {
	Default string
	Id      string
	Type    string
	UUID    string
}

func NewIdentifierGenerator(formatString string) *IdentifierGenerator {
	return &IdentifierGenerator{formatString: formatString}
}

func (generator *IdentifierGenerator) NewIdentifier(attributes IdentifierAttributes) string {
	replacer := strings.NewReplacer("%{uuid}", attributes.UUID, "%{id}", attributes.Id, "%{type}", attributes.Type, "%{default}", attributes.Default)
	return replacer.Replace(generator.formatString)
}
