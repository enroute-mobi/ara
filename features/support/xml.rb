require 'rexml/document'
require 'rexml/xpath'

def normalized_xml(xml)
  "".tap do |output|
    REXML::Document.new(xml).write output: output, indent: 2
  end
end
