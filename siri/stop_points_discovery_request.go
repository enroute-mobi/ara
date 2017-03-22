package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopDiscoveryRequest struct {
	RequestXMLStructure
}

const stopDiscoveryRequestTemplate = `<ns7:StopPointsDiscovery xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <Request>
              <ns2:RequestTimestamp>2017-03-03T11:28:00.359Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>STIF</ns2:RequestorRef>
            </Request>
            <RequestExtension />
						</ns7:StopPointsDiscovery>`

func NewXMLStopDiscoveryRequest(node xml.Node) *XMLStopDiscoveryRequest {
	xmlStopDiscoveryRequest := &XMLStopDiscoveryRequest{}
	xmlStopDiscoveryRequest.node = NewXMLNode(node)
	return xmlStopDiscoveryRequest
}

func NewXMLStopDiscoveryRequestFromContent(content []byte) (*XMLStopDiscoveryRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopDiscoveryRequest(doc.Root().XmlNode)
	return request, nil
}
