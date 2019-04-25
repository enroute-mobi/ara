package siri

import (
	"bytes"
	"text/template"
	"time"
)

type SIRINotifyStopMonitoring struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	ResponseTimestamp         time.Time

	Deliveries []*SIRINotifyStopMonitoringDelivery
}

type SIRINotifyStopMonitoringDelivery struct {
	MonitoringRef          string
	RequestMessageRef      string
	SubscriberRef          string
	SubscriptionIdentifier string
	ResponseTimestamp      time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	MonitoredStopVisits []*SIRIMonitoredStopVisit
	CancelledStopVisits []*SIRICancelledStopVisit
}

const stopMonitoringNotifyTemplate = `<sw:NotifyStopMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{ .ProducerRef }}</siri:ProducerRef>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification>{{ range .Deliveries }}
		{{ .BuildNotifyStopMonitoringDeliveryXML }}{{ end }}
	</Notification>
	<NotifyExtension />
</sw:NotifyStopMonitoring>`

const notifyStopMonitoringDeliveryTemplate = `<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
			<siri:SubscriberRef>{{ .SubscriberRef }}</siri:SubscriberRef>
			<siri:SubscriptionRef>{{ .SubscriptionIdentifier }}</siri:SubscriptionRef>{{ if .MonitoringRef }}
			<siri:MonitoringRef>{{ .MonitoringRef }}</siri:MonitoringRef>{{ end }}
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{ .ErrorNumber }}">{{ else }}
				<siri:{{ .ErrorType }}>{{ end }}
					<siri:ErrorText>{{ .ErrorText }}</siri:ErrorText>
				</siri:{{ .ErrorType }}>
			</siri:ErrorCondition>{{ else }}{{ range .MonitoredStopVisits }}
			{{ .BuildMonitoredStopVisitXML }}{{ end }}{{ range .CancelledStopVisits }}
			{{ .BuildCancelledStopVisitXML }}{{ end }}{{ end }}
		</siri:StopMonitoringDelivery>`

func (notify *SIRINotifyStopMonitoring) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("stopMonitoringNotify").Parse(stopMonitoringNotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRINotifyStopMonitoringDelivery) BuildNotifyStopMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	var stopMonitoringDelivery = template.Must(template.New("notifyStopMonitoringDelivery").Parse(notifyStopMonitoringDeliveryTemplate))
	if err := stopMonitoringDelivery.Execute(&buffer, delivery); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
