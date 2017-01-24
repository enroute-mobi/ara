require 'rexml/document'
require 'rexml/xpath'

def siri_path(attributes = {})
  url_for(attributes.merge(path: "siri"))
end

Given(/^we send a checkstatus request for referential "([^"]+)"$/) do |referential|
  xmlBody = File.read("features/testdata/checkstatus-soap-request.xml")
  response = RestClient.post siri_path(referential: referential), xmlBody, {content_type: :xml}
  @last_siri_response = response.body
end

Then(/^we should recieve a positive checkstatus response$/) do
  xmlBody = @last_siri_response
  doc = REXML::Document.new xmlBody
  status = REXML::XPath.first(doc, "//*[local-name()='Status']")
  expect(status.text).to eq("true")
end
