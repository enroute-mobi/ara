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

	objectidkind = attributes["ObjectIDs"].keys.first
  objectid_value = attributes["ObjectIDs"][objectidkind]

  expectedAttr = responseArray.find{|a| a["ObjectIDs"][objectidkind] == objectid_value }

  expect(expectedAttr).not_to be_nil
end

Then(/^a VehicleJourney "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get vehicle_journeys_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedVihicleJourney = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  if condition.nil?
    expect(expectedLine).not_to be_nil
  else
    expect(expectedLine).to be_nil
  end
end