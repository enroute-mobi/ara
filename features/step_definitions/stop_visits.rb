def stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit'))
end

def stop_visit_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit', id: id))
end


Given(/^a StopVisit exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_visit|
  RestClient.post stop_visits_path(referential: referential), model_attributes(stop_visit).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

When(/^a StopVisit is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  if referential.nil?
    step "a StopVisit exists with the following attributes:", stopArea
  else
    step "a StopVisit exists in Referential \"#{referential}\" with the following attributes:", stopArea
  end
end

When(/^the StopVisit "([^"]*)" is edited with the following attributes:$/) do |identifier, attributes|
  RestClient.put stop_visit_path(identifier), model_attributes(attributes).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

Then(/^the StopVisit "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  # puts RestClient.get stop_visits_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response = RestClient.get stop_visit_path(identifier), {content_type: :json, :Authorization => "Token token=#{$token}"}
  stopVisitAttributes = api_attributes(response.body)
  expect(stopVisitAttributes).to include(model_attributes(attributes))
end

Then(/^one StopVisit has the following attributes:$/) do |attributes|
  # puts RestClient.get stop_visits_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response = RestClient.get stop_visits_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end


Then(/^a StopVisit "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
  response = RestClient.get(stop_visit_path("#{kind}:#{value}", referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}){|response, request, result| response }

  if condition.nil?
    expect(response.code).to eq(200)
  else
    expect(response.code).to eq(404)
    expect(response.body).to include("Stop visit not found: #{kind}:#{value}")
  end
end
