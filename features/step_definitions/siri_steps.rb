# coding: utf-8
require 'rexml/document'
require 'rexml/xpath'

def siri_path(attributes = {})
  attributes = {
    referential: 'test'
  }.merge(attributes.delete_if { |k,v| v.nil? })

  url_for(attributes.merge(path: "siri"))
end

When(/^I send this SIRI request(?: to the Referential "([^"]*)")?$/) do |referential, request|
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
