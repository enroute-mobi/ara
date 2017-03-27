def stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit'))
end

def stop_visit_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit', id: id))
end

Given(/^a StopVisit exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_visit|
  RestClient.post stop_visits_path(referential: referential), model_attributes(stop_visit).to_json, {content_type: :json}
end

When(/^a StopVisit is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  if referential.nil?
    step "a StopVisit exists with the following attributes:", stopArea
  else
    step "a StopVisit exists in Referential \"#{referential}\" with the following attributes:", stopArea
  end
end

Then(/^the StopVisit "([^"]*)" has the following attributes:$/) do |identifier, attributes|
	response = RestClient.get stop_visit_path(identifier)
	stopVisitAttributes = api_attributes(response.body)
  expect(stopVisitAttributes).to include(model_attributes(attributes))
end

Then(/^one StopVisit has the following attributes:$/) do |attributes|
  response = RestClient.get stop_visits_path
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end


Then(/^a StopVisit "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get stop_visits_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedStopVisit = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  if condition.nil?
    expect(expectedStopVisit).not_to be_nil
  else
    expect(expectedStopVisit).to be_nil
  end
end
