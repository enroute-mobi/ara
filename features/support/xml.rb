require 'rexml/document'
require 'rexml/xpath'

def normalized_xml(xml)
  "".tap do |output|
    REXML::Document.new(xml).write output: output, indent: 2
  end
end

module XML
  class Document

    def initialize(content)
      @document = REXML::Document.new(content)
    end

    def normalize
      "".tap do |output|
        @document.write output: output, indent: 2
      end
    end

    def values(xpaths = [])
      {}.tap do |values|
        xpaths.each do |xpath|
          node = REXML::XPath.first(@document, xpath, { "siri" => "http://www.siri.org.uk/siri" })
          xml_value = node.text if node
          values[xpath] = xml_value
        end
      end
    end
  end
end
