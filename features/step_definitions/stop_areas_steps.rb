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
	responseArray = JSON.parse(response.body)

	stopAreaHash = model_attributes(stopArea)
	objectidkind = stopAreaHash["ObjectIds"].keys.first
	objectid_value = stopAreaHash["ObjectIds"][objectidkind]

	expectedName = responseArray.find{|a| a["Name"] == stopAreaHash["Name"]}
	expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

	expect(expectedName).not_to be_nil
	expect(expectedAttr).not_to be_nil
end


Then(/^a StopArea "([^"]+)":"([^"]+)" should exist$/) do |kind, objectid|
	response = RestClient.get url
	responseArray = JSON.parse(response.body)
	expectation = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == kind && o["Value"] == objectid }}
	puts expectation
	expect(expectation).to be_truthy
end

Given(/^a StopArea exists with the following attributes :$/) do |stopArea|
 	creation = RestClient.post url, model_attributes(stopArea).to_json, {content_type: :json, accept: :json}

 	expect(creation).not_to be_nil
end

When(/^the StopArea "([^"]+)":"([^"]+)" is destroy :$/) do |kind, objectid|
	response = RestClient.get url
	puts responseArray = JSON.parse(response.body)
	kInd = kind
	id = objectid

	RestClient.delete "#{url}/#{kInd}/#{id}"

end