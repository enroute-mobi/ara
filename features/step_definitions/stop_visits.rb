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
  responseArray = JSON.parse(response.body)

  attributes = model_attributes(attributes)

  objectidkind = attributes["ObjectIDs"].keys.first
  objectid_value = attributes["ObjectIDs"][objectidkind]


  expectedName = responseArray.find{|a| a["Name"] == attributes["Name"]}
  expectedAttr = responseArray.find{|a| a["ObjectIDs"][objectidkind] == objectid_value }

  expect(expectedName).not_to be_nil
  expect(expectedAttr).not_to be_nil
end
