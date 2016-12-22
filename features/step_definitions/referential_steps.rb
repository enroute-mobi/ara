require 'rest-client'
require 'json'

url = "http://localhost:8081/_referentials"

When(/^a Referential "([^"]+)" is created$/) do |referential|
	RestClient.post url,"{\"slug\":\"#{referential}\"}", {content_type: :json, accept: :json}
end

Then(/^one Referential "([^"]+)" should exist$/) do |referential|
	response = RestClient.get url
	responseHash = JSON.parse(response.body)

	expect(responseHash.find{|a| a["Slug"] == referential}).not_to be_nil
end

Given(/^a Referential "([^"]+)" exists$/) do |referential|
	step "a Referential \"#{referential}\" is created"
end

When(/^the Referential "([^"]+)" is destroyed$/) do |referential|
	response = RestClient.get url
	responseHash = JSON.parse(response.body)

	id = responseHash.find{|a| a["Slug"] == referential}["Id"]
	RestClient.delete "#{url}/#{id}"
end

Then(/^a Referential "([^"]+)" should not exists$/) do |referential|
	response = RestClient.get url
	responseHash = response.body == "null" ? [] : JSON.parse(response.body)

	expect(responseHash.find{|a| a["Slug"] == referential}).to be_nil
end

