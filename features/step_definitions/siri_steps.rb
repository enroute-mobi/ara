# coding: utf-8
require 'rexml/document'
require 'rexml/xpath'

def siri_path(attributes = {})
  url_for(attributes.merge(path: "siri"))
end

When(/^we send a checkstatus request for referential "([^"]+)"$/) do |referential|
  xmlBody = File.read("features/testdata/checkstatus-soap-request.xml")
  response = RestClient.post siri_path(referential: referential), xmlBody, {content_type: :xml}
  @last_siri_response = response.body
end

Then(/^we should receive a positive checkstatus response$/) do
  xmlBody = @last_siri_response
  doc = REXML::Document.new xmlBody
  status = REXML::XPath.first(doc, "//*[local-name()='Status']")
  expect(status.text).to eq("true")
end

When(/^I send SIRI request to the referential "([^"]*)"$/) do |referential, request|
  response = RestClient.post siri_path(referential: referential), request, {content_type: :xml}
  @last_siri_response = response.body
end

def normalized_xml(xml)
  "".tap do |output|
    REXML::Document.new(xml).write output: output, indent: 2
  end
end

Then(/^I should receive this SIRI reponse$/) do |expected_xml|
  expect(normalized_xml(@last_siri_response)).to eq(normalized_xml(expected_xml))
end
