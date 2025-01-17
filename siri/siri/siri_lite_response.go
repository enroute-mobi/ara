package siri

import (
	"encoding/json"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type SiriLiteResponse struct {
	Siri *SiriLiteResponseSubstructure
}

type SiriLiteResponseSubstructure struct {
	ServiceDelivery *ServiceDelivery
}

type ServiceDelivery struct {
	ResponseTimestamp         time.Time                      `json:",omitempty"`
	ProducerRef               string                         `json:",omitempty"`
	ResponseMessageIdentifier string                         `json:",omitempty"`
	RequestMessageRef         string                         `json:",omitempty"`
	VehicleMonitoringDelivery *SIRIVehicleMonitoringDelivery `json:",omitempty"`
	// StopMonitoringDelivery *StopMonitoringDelivery `json:",omitempty"`
	// ...
}

type ErrorCondition struct {
	ErrorType   string
	ErrorNumber int
	ErrorText   string
}

func (ec *ErrorCondition) MarshalJson() ([]byte, error) {
	aux := make(map[string]map[string]string)
	aux[ec.ErrorType] = make(map[string]string)
	if ec.ErrorType == siri_attributes.OtherError {
		aux[ec.ErrorType]["number"] = strconv.Itoa(ec.ErrorNumber)
	}
	aux[ec.ErrorType]["ErrorText"] = ec.ErrorText

	return json.Marshal(aux)
}

func NewSiriLiteResponse() *SiriLiteResponse {
	return &SiriLiteResponse{
		Siri: &SiriLiteResponseSubstructure{ServiceDelivery: &ServiceDelivery{}},
	}
}
