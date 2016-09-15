package siri

import (
	"fmt"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/jbowtie/gokogiri/xpath"
)

type XMLCheckStatusRequest struct {
	content []byte
}

func NewXMLCheckStatusRequest(content []byte) *XMLCheckStatusRequest {
	return &XMLCheckStatusRequest{content: content}
}

func (request *XMLCheckStatusRequest) Document() *xml.XmlDocument {
	doc, _ := gokogiri.ParseXml(request.content)
	// defer doc.Free()
	return doc
}

func (request *XMLCheckStatusRequest) RequestorRef() string {
	path := xpath.Compile("//*[local-name()='RequestorRef']")
	nodes, _ := request.Document().Root().Search(path)
	fmt.Println(nodes)
	return nodes[0].Content()
}
