def vehicle_journeys_path(attributes = {})
	url_for_model(attributes.merge(resource: 'vehicle_journey'))
end

Given(/^a VehicleJourney exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle_journey|
  response = RestClient.post vehicle_journeys_path(referential: referential), model_attributes(vehicle_journey).to_json, {content_type: :json}
  # puts response.body
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
	response = RestClient.get vehicle_journeys_path
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end

Then(/^a VehicleJourney "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get vehicle_journeys_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedVehicleJourney = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  if condition.nil?
    expect(expectedVehicleJourney).not_to be_nil
  else
    expect(expectedVehicleJourney).to be_nil
  end
end
