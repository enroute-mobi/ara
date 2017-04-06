def stop_areas_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_area'))
end

def stop_area_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'stop_area', id: id))
end

Given(/^a StopArea exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  RestClient.post stop_areas_path(referential: referential), model_attributes(stopArea).to_json, {content_type: :json}
end

When(/^a StopArea is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  if referential.nil?
    step "a StopArea exists with the following attributes:", stopArea
  else
    step "a StopArea exists in Referential \"#{referential}\" with the following attributes:", stopArea
  end
end

When(/^the StopArea "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, value, referential|
  response = RestClient.get stop_area_path("#{kind}:#{value}", referential: referential)
  expectedStopArea = JSON.parse(response.body)

  RestClient.delete stop_area_path expectedStopArea["Id"]
end

Then(/^one StopArea(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get stop_areas_path(referential: referential)
  response_array = api_attributes(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end

Then(/^a StopArea "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
 response = RestClient.get(stop_area_path("#{kind}:#{value}", referential: referential)){|response, request, result| response }
  
  if condition.nil?
    expect(response.code).to eq(200)
  else
    expect(response.code).to eq(500)
    expect(response.body).to include("Stop area not found: #{kind}:#{value}")
  end
end
