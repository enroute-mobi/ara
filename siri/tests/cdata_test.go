package siri_tests

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

func Test_CData(t *testing.T) {
	input := `<DATA>
    <NAME><![CDATA[FIRSTNAME LASTNAME MIDDLENAME]]></NAME>
    <NUM>3731</NUM>
    <person_type>4</person_type>
    <birth_date><![CDATA[01.11.1992]]></birth_date>
    <DESCRIPTION><![CDATA[DESCRIPTION]]></DESCRIPTION>
</DATA>`
	doc, err := gokogiri.ParseXml([]byte(input))
	if err != nil {
		t.Fatal("Can't parse input")
	}
	n, err := doc.Root().Search("//DATA/NAME")
	if err != nil {
		t.Fatal("Can't search")
	}
	// var e *xml.ElementNode
	// e.Search()
	logger.Log.Printf("%+v", n[0])
	var node xml.Node
	node = n[0]
	logger.Log.Printf("Node: %+v", node)
	logger.Log.Printf("Node string: %+v", node.String())
	logger.Log.Printf("Node content: %+v", node.Content())
	logger.Log.Printf("Node type: %v", node.NodeType())
	logger.Log.Printf("Node FirstChild: %v", node.FirstChild())
	logger.Log.Printf("Node FirstChild Type: %v", node.FirstChild().NodeType())
	logger.Log.Printf("Node FirstChild Content: %v", node.FirstChild().Content())
	logger.Log.Printf("Node FirstChild String: %v", node.FirstChild().String())

	t.Fatal("pouet")
	// => "FIRSTNAME LASTNAME MIDDLENAME "

	// doc.xpath('//DATA').each do |terr|
	//   puts "\nName: "+terr.xpath('NAME').text
	// end

	// => Name: FIRSTNAME LASTNAME MIDDLENAME
}
