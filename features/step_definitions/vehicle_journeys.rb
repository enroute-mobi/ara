def vehicle_journey_path(attributes = {})
	url_for_model(attributes.merge(resource: 'vehicle_journey'))
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
	response = RestClient.get vehicle_journey_path
	responseArray = JSON.parse(response.body)

	attributes = model_attributes(attributes)

	objectidkind = attributes["ObjectIds"].keys.first
  objectid_value = attributes["ObjectIds"][objectidkind]

  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedAttr).not_to be_nil
end
