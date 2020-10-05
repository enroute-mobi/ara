package core

import (
	"regexp"
	"strings"

	"bitbucket.org/enroute-mobi/edwig/model"
)

var defaultIdentifierGenerators = map[string]string{
	"message_identifier":             "%{uuid}",
	"response_message_identifier":    "%{uuid}",
	"data_frame_identifier":          "%{id}",
	"reference_identifier":           "%{type}:%{default}",
	"reference_stop_area_identifier": "%{default}",
	"subscription_identifier":        "%{id}",
}

type IdentifierGenerator struct {
	model.UUIDConsumer

	formatString string
	pattern      string
	replacement  string
}

type IdentifierAttributes struct {
	Default  string
	Id       string
	ObjectId string
	Type     string
}

func NewIdentifierGenerator(formatString string) *IdentifierGenerator {
	objectidSubstitutionPattern := regexp.MustCompile(`%{objectid//([^/]+)/([^/]*)}`)
	matches := objectidSubstitutionPattern.FindStringSubmatch(formatString)

	var pattern, replacement string

	if len(matches) == 3 {
		pattern = matches[1]
		replacement = matches[2]

		formatString = objectidSubstitutionPattern.ReplaceAllString(formatString, `%{objectid}`)
	}

	return &IdentifierGenerator{formatString: formatString, pattern: pattern, replacement: replacement}
}

func NewIdentifierGeneratorWithUUID(formatString string, uuidGenerator model.UUIDConsumer) *IdentifierGenerator {
	generator := NewIdentifierGenerator(formatString)
	generator.UUIDConsumer = uuidGenerator
	return generator
}

func (generator *IdentifierGenerator) NewIdentifier(attributes IdentifierAttributes) string {
	objectidValue := attributes.ObjectId

	if len(generator.pattern) > 0 {
		objectidValue = strings.ReplaceAll(objectidValue, generator.pattern, generator.replacement)
	}

	replacer := strings.NewReplacer("%{id}", attributes.Id, "%{type}", attributes.Type, "%{default}", attributes.Default, "%{objectid}", objectidValue)
	return generator.handleuuids(replacer.Replace(generator.formatString))
}

func (generator *IdentifierGenerator) NewMessageIdentifier() string {
	return generator.handleuuids(generator.formatString)
}

func (generator *IdentifierGenerator) handleuuids(s string) string {
	re := regexp.MustCompile("%{uuid}")
	return re.ReplaceAllStringFunc(s, func(string) string { return generator.NewUUID() })
}
