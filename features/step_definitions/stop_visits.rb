def stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit'))
end

Then(/^one StopVisit has the following attributes:$/) do |attributes|
	response = RestClient.get stop_visits_path
	responseArray = JSON.parse(response.body)

	attributes = model_attributes(attributes)

	objectidkind = attributes["ObjectIds"].keys.first
  objectid_value = attributes["ObjectIds"][objectidkind]

	expectedName = responseArray.find{|a| a["Name"] == attributes["Name"]}
  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedName).not_to be_nil
	expect(expectedAttr).not_to be_nil
end
