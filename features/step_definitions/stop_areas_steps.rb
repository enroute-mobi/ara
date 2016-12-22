require 'rest-client'
require 'json'

url = "http://localhost:8081/test/stop_areas"

def model_attributes(table)
	attributes = table.rows_hash
	if attributes["ObjectIds"]
		attributes["ObjectIds"] = JSON.parse("{" + attributes["ObjectIds"] + "}")
	end
	attributes
end


When(/^a StopArea is created with the following attributes :$/) do |stopArea|
	RestClient.post url, model_attributes(stopArea).to_json, {content_type: :json, accept: :json}
end

Then(/^one StopArea has the following attributes:$/) do |stopArea|
	response = RestClient.get url
	puts responseHash = JSON.parse(response.body)
	expectation = responseHash.find{|a| puts a[model_attributes(stopArea)] == a[stopArea]}

	expect(expectation).to be_truthy
end


Then(/^a StopArea "([^"]+)":"([^"]+)" should exist$/) do |kind, objectid|
	response = RestClient.get url
	responseHash = JSON.parse(response.body)
	expect(responseHash.find{|a| a[kind] = objectid}).to be_truthy
	puts responseHash.find{|a| a[kind] = objectid}
end

# Given(/^a StopArea exists with the following attributes :$/) do |stopArea|
# 	RestClient.post url, model_attributes(stopArea).to_json, {content_type: :json, accept: :json}
# 	response = RestClient.get url
# 	responseHash = JSON.parse(response.body)
# 	expect(responseHash.find{|a| a["ObjectIds"] = stopArea}).to be_truthy
# end