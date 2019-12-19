def vehicles_path(attributes = {})
  url_for_model(attributes.merge(resource: 'vehicle'))
end

def vehicle_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'vehicle', id: id))
end


Given(/^a Vehicle exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle|
  RestClient.post vehicles_path(referential: referential), model_attributes(vehicle).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

When(/^a Vehicle is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle|
  if referential.nil?
    step "a Vehicle exists with the following attributes:", vehicle
  else
    step "a Vehicle exists in Referential \"#{referential}\" with the following attributes:", vehicle
  end
end

When(/^the Vehicle "([^"]*)" is edited with the following attributes:$/) do |identifier, attributes|
  RestClient.put vehicle_path(identifier), model_attributes(attributes).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
  # puts RestClient.get vehicles_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

Then(/^the Vehicle "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  # puts RestClient.get vehicles_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response = RestClient.get vehicle_path(identifier), {content_type: :json, :Authorization => "Token token=#{$token}"}
  vehicleAttributes = api_attributes(response.body)
  expect(vehicleAttributes).to include(model_attributes(attributes))
end

Then(/^one Vehicle has the following attributes:$/) do |attributes|
  # puts RestClient.get vehicles_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response = RestClient.get vehicles_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end


Then(/^a Vehicle "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
  response = RestClient.get(vehicle_path("#{kind}:#{value}", referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"})
  if condition.nil?
    expect(response.code).to eq(200)
  else
    expect(response.code).to eq(404)
    expect(response.body).to include("Stop visit not found: #{kind}:#{value}")
  end
end
