package sxml

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSiriError struct {
	XMLStructure

	e string
}

func NewXMLSiriError(node xml.Node) *XMLSiriError {
	se := &XMLSiriError{}
	se.node = NewXMLNode(node)
	return se
}

func (se *XMLSiriError) Error() string {
	if se.e == "" {
		se.e = fmt.Sprintf("%v: %v", se.findStringChildContent(siri_attributes.FaultCode), se.findStringChildContent(siri_attributes.FaultString))
	}
	return se.e
}
