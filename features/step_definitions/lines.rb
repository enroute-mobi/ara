def lines_path(attributes = {})
  url_for_model(attributes.merge(resource: 'line'))
end

def line_path(id, attributes = {})
  path = url_for_model(attributes.merge(resource: 'line', id: id))
end

Given(/^a Line exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line|
  response = RestClient.post lines_path(referential: referential), model_attributes(line).to_json, {content_type: :json, :Authorization => "Token token=#{$token}" }
  # puts response.body
end

When(/^a Line is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line|
  if referential.nil?
    step "a Line exists with the following attributes:", line
  else
    step "a Line exists in Referential \"#{referential}\" with the following attributes:", line
  end
end

When(/^the Line "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, objectid, referential|
  response = RestClient.get lines_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  responseArray = JSON.parse(response.body)
  expectedLine = responseArray.find{|a| a["ObjectIDs"][kind] == objectid }

  RestClient.delete line_path(expectedLine["Id"]), {:Authorization => "Token token=#{$token}"}
end

Then(/^one Line(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get lines_path(referential: referential), {:Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end


Then(/^a Line "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
  response = RestClient.get(line_path("#{kind}:#{value}" ,referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"} ){|response, request, result| response }

  if condition.nil?
    expect(response.code).to eq(200)
  else
    expect(response.code).to eq(404)
    expect(response.body).to include("Line not found: #{kind}:#{value}")
  end
end
