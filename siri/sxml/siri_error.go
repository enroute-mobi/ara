package sxml

import (
	"fmt"

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
		se.e = fmt.Sprintf("%v: %v", se.findStringChildContent("faultcode"), se.findStringChildContent("faultstring"))
	}
	return se.e
}
