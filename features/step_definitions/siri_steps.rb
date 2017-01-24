require 'rexml/document'
require 'rexml/xpath'

working_directory = "features/testdata"

def siri_path(attributes = {})
  url_for(attributes.merge(path: "siri"))
end

Given(/^we send a checkstatus request for referential "([^"]+)"$/) do |referential|
  system "rm #{working_directory}/response.xml"
  xmlBody = File.read("#{working_directory}/checkstatus-soap-request.xml")
  response = RestClient.post siri_path(referential: referential), xmlBody, {content_type: :xml}
  File.write("#{working_directory}/response.xml", response.body)
end

Then(/^we should recieve a positive checkstatus response$/) do
  xmlBody = File.read("#{working_directory}/response.xml")
  doc = REXML::Document.new xmlBody
  status = REXML::XPath.first(doc, "//*[local-name()='Status']")

  expect(status.text).to eq("true")
end
