def vehicle_journeys_path(attributes = {})
	url_for_model(attributes.merge(resource: 'vehicle_journey'))
end

def vehicle_journey_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'vehicle_journey', id: id))
end

Given(/^a VehicleJourney exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle_journey|
  response = RestClient.post vehicle_journeys_path(referential: referential), model_attributes(vehicle_journey).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
  debug response.body
end

When(/^the VehicleJourney "([^"]*)" is edited with the following attributes:$/) do |identifier, attributes|
  RestClient.put vehicle_journey_path(identifier), model_attributes(attributes).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
  # puts RestClient.get vehicles_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

Then(/^the VehicleJourney "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  response = RestClient.get vehicle_journey_path(identifier), {content_type: :json, :Authorization => "Token token=#{$token}"}
  vehicleJourneyAttributes = api_attributes(response.body)
  expect(vehicleJourneyAttributes).to include(model_attributes(attributes))
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
	response = RestClient.get vehicle_journeys_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end

Then(/^a VehicleJourney "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
  response = RestClient.get(vehicle_journey_path("#{kind}:#{value}", referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}){|response, request, result| response }

  if condition.nil?
    expect(response.code).to eq(200)
  else
    expect(response.code).to eq(404)
    expect(response.body).to include("Vehicle journey not found: #{kind}:#{value}")
  end
end
