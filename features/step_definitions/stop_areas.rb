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

When(/^the StopArea "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, objectid, referential|
  response = RestClient.get stop_areas_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedStopArea = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  RestClient.delete stop_area_path expectedStopArea["Id"]
end

Then(/^one StopArea(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get stop_areas_path(referential: referential)
  response_array = api_attributes(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end

Then(/^a StopArea "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get stop_areas_path(referential: referential)
  stopAreas = api_attributes(response.body)
  expectedStopArea = stopAreas.find{|a| a["ObjectIDs"][kind] == objectid }

  if condition.nil?
    expect(expectedStopArea).not_to be_nil
  else
    expect(expectedStopArea).to be_nil
  end
end
