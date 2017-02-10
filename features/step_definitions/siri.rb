def siri_path(attributes = {})
  attributes = {
    referential: 'test'
  }.merge(attributes.delete_if { |k,v| v.nil? })

  url_for(attributes.merge(path: "siri"))
end

Given(/^a SIRI server waits (GetStopMonitoring) request on "([^"]*)" to respond with$/) do |message_type, url, response|
  (@the_siri_server = SIRIServer.create(url)).expect_request(message_type, response).start
end

When(/^the SIRI server has received a (GetStopMonitoring) request$/) do |message_type|
  @the_siri_server.wait_request message_type
end

When(/^I send this SIRI request(?: to the Referential "([^"]*)")?$/) do |referential, request|
  response = RestClient.post siri_path(referential: referential), request, {content_type: :xml}
  @last_siri_response = response.body
end

Then(/^I should receive this SIRI reponse$/) do |expected_xml|
  expect(normalized_xml(@last_siri_response)).to eq(normalized_xml(expected_xml))
end
