package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func CancelSubscription(subId, kind string, connector Connector) {
	message := &audit.BigQueryMessage{
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.Partner().Slug()),
		Status:    "OK",
	}
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.Partner().ProducerRef(),
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
	}
	logSIRIDeleteSubscriptionRequest(message, request, kind, connector.Partner().SIRIEnvelopeType())

	startTime := connector.Clock().Now()
	response, err := connector.Partner().SIRIClient().DeleteSubscription(request)

	responseTime := connector.Clock().Since(startTime)
	message.ProcessingTime = responseTime.Seconds()

	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : %v", subId, err.Error())
		e := fmt.Sprintf("Error during DeleteSubscription: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}
	logXMLDeleteSubscriptionResponse(message, response)
}

func logSIRIDeleteSubscriptionRequest(message *audit.BigQueryMessage, request *siri.SIRIDeleteSubscriptionRequest, subType, envelopeType string) {
	message.Type = "DeleteSubscriptionRequest"
	message.RequestIdentifier = request.MessageIdentifier
	message.SubscriptionIdentifiers = []string{request.SubscriptionRef}

	xml, err := request.BuildXML(envelopeType)
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLDeleteSubscriptionResponse(message *audit.BigQueryMessage, response *sxml.XMLDeleteSubscriptionResponse) {
	var i int
	for _, responseStatus := range response.ResponseStatus() {
		if !responseStatus.Status() {
			i++
		}
	}

	if i > 0 {
		message.Status = "Error"
		message.ErrorDetails = fmt.Sprintf("%d ResponseStatus returned false", i)
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
