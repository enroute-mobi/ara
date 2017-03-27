def lines_path(attributes = {})
  url_for_model(attributes.merge(resource: 'line'))
end

def line_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'line', id: id))
end

Given(/^a Line exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line|
  RestClient.post lines_path(referential: referential), model_attributes(line).to_json, {content_type: :json}
end

When(/^a Line is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line|
  if referential.nil?
    step "a Line exists with the following attributes:", line
  else
    step "a Line exists in Referential \"#{referential}\" with the following attributes:", line
  end
end

When(/^the Line "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, objectid, referential|
  response = RestClient.get lines_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedLine = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  RestClient.delete line_path expectedLine["Id"]
end

Then(/^one Line(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get lines_path(referential: referential)
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end


Then(/^a Line "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get lines_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedLine = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  if condition.nil?
    expect(expectedLine).not_to be_nil
  else
    expect(expectedLine).to be_nil
  end
end
