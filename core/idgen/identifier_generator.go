package idgen

import (
	"regexp"
	"strings"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	MessageSetting         = "generators.message_identifier"
	ResponseMessageSetting = "generators.response_message_identifier"
	DataFrameSetting       = "generators.data_frame_identifier"
	ReferenceSetting       = "generators.reference_identifier"
	StopAreaSetting        = "generators.reference_stop_area_identifier"
	VehicleJourneySetting  = "generators.reference_vehicle_journey_identifier"
	SubscriptionSetting    = "generators.subscription_identifier"

	Message         = "Message"
	ResponseMessage = "ResponseMessage"
	DataFrame       = "DataFrame"
	Reference       = "Reference"
	StopArea        = "StopArea"
	VehicleJourney  = "VehicleJourney"
	Subscription    = "Subscription"
)

var generatorSettings = map[string]string{
	MessageSetting:         Message,
	ResponseMessageSetting: ResponseMessage,
	DataFrameSetting:       DataFrame,
	ReferenceSetting:       Reference,
	StopAreaSetting:        StopArea,
	VehicleJourneySetting:  VehicleJourney,
	SubscriptionSetting:    Subscription,
}

var re = regexp.MustCompile("%{uuid}")

var defaultGenerators = map[string]string{
	Message:         "%{uuid}",
	ResponseMessage: "%{uuid}",
	DataFrame:       "%{id}",
	Reference:       "%{type}:%{id}",
	StopArea:        "%{id}",
	Subscription:    "%{id}",
}

type IdentifierGenerator struct {
	uuidGenerator uuid.UUIDGenerator
	formatStrings map[string]string
}

type IdentifierAttributes struct {
	Type string
	Id   string
}

func NewIdentifierGenerator(settings map[string]string, uuidGenerator uuid.UUIDGenerator) IdentifierGenerator {
	var formatStrings = make(map[string]string, len(generatorSettings))

	/* If the setting is defined we use it
	if not we check if there's a default one
	if not, we don't do anything and Reference will be used
	*/
	for setting, kind := range generatorSettings {
		format, ok := settings[setting]
		if !ok {
			format, ok = defaultGenerators[kind]
			if !ok {
				continue
			}
		}
		formatStrings[kind] = format
	}

	return IdentifierGenerator{
		uuidGenerator: uuidGenerator,
		formatStrings: formatStrings,
	}
}

// We use attributes.Type to find the appropriate generator
// We can use the kind optionnal argument if Type isn't clear enough like for DatedVehicleJourney
func (idg *IdentifierGenerator) NewIdentifier(attributes IdentifierAttributes, kind ...string) string {
	// default and code are legacy values, keep them for now for a smoother transition
	replacer := strings.NewReplacer("%{id}", attributes.Id, "%{type}", attributes.Type, "%{default}", attributes.Id, "%{code}", attributes.Id)
	return idg.handleuuids(replacer.Replace(idg.formatString(kind, attributes)))
}

// To avoid the replacer
func (idg *IdentifierGenerator) NewMessageIdentifier() string {
	return idg.handleuuids(idg.formatStrings[Message])
}

func (idg *IdentifierGenerator) NewResponseMessageIdentifier() string {
	return idg.handleuuids(idg.formatStrings[ResponseMessage])
}

func (idg *IdentifierGenerator) handleuuids(s string) string {
	return re.ReplaceAllStringFunc(s, func(string) string { return idg.uuidGenerator.NewUUID() })
}

// We check if we know a generator for the identifier, otherwise we consider it's a Reference
func (idg *IdentifierGenerator) formatString(kind []string, attributes IdentifierAttributes) string {
	if s, ok := idg.formatStrings[processKind(kind, attributes)]; ok {
		return s
	}
	return idg.formatStrings[Reference]
}

/* We check in order :
* If a kind is specified
* If a Type is specified
Otherwise we return 'Reference'
*/
func processKind(k []string, a IdentifierAttributes) string {
	if len(k) != 0 {
		return k[0]
	}

	if a.Type != "" {
		return a.Type
	}

	return Reference
}

// Test method
func (idg *IdentifierGenerator) FormatString(kind string) string {
	return idg.formatStrings[kind]
}
