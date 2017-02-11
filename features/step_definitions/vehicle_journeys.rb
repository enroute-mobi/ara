def vehicle_journeys_path(attributes = {})
	url_for_model(attributes.merge(resource: 'vehicle_journey'))
end

Given(/^a VehicleJourney exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle_journey|
  response = RestClient.post vehicle_journeys_path(referential: referential), model_attributes(vehicle_journey).to_json, {content_type: :json}
  # puts response.body
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
	response = RestClient.get vehicle_journeys_path
	responseArray = JSON.parse(response.body)

	attributes = model_attributes(attributes)

	objectidkind = attributes["ObjectIds"].keys.first
  objectid_value = attributes["ObjectIds"][objectidkind]

  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedAttr).not_to be_nil
end