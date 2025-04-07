def referentials_path
  url_for(path: "_referentials")
end

def referential_path(id)
  url_for(path: "_referentials/#{id}")
end

Given('a Referential {string} exists with the following attributes:') do |referential, table|
  user_attributes = table.rows_hash

  attributes = {
    Slug: referential,
    Tokens: user_attributes.fetch('Tokens', $token).split(','),
    ImportTokens: user_attributes.fetch('Import Tokens', '').split(',')
  }

  RestClient.post referentials_path, attributes.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
end

Given(/^a Referential "([^"]+)" exists$/) do |referential|
  attributes = {
    slug: referential,
    tokens: [$token]
  }
  RestClient.post referentials_path, attributes.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
end

Given(/^a Referential "([^"]*)" exists with the following settings:$/) do |referential, settings|
  attributes = {
    slug: referential,
    tokens: [$token]
  }
  attributes[:settings] = settings.rows_hash if settings
  RestClient.post referentials_path, attributes.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
end

When(/^a Referential "([^"]+)" is created$/) do |referential|
  step "a Referential \"#{referential}\" exists"
end

When(/^the Referential "([^"]+)" is destroyed$/) do |referential|
  response = RestClient.get referentials_path, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
  responseHash = JSON.parse(response.body)

  id = responseHash.find{|a| a["Slug"] == referential}["Id"]
  RestClient.delete referential_path(id), {:Authorization => "Token token=#{$adminToken}"}
end

Then(/^a Referential "([^"]+)" should (not )?exist$/) do |referential, condition|
  response = RestClient.get referentials_path, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
  responseHash = JSON.parse(response.body)

  if condition.nil?
    expect(responseHash.find{|a| a["Slug"] == referential}).not_to be_nil
  else
    expect(responseHash.find{|a| a["Slug"] == referential}).to be_nil
  end
end

Then(/^one Referential has the following attributes:$/) do |attributes|
  response = RestClient.get referentials_path, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
  response_array = api_attributes(response.body)

  parsed_attributes = model_attributes(attributes)
  found_value = response_array.find{|a| a["Id"] == parsed_attributes["Id"]}

  expect(found_value).not_to be_nil

  expect(found_value).to include(parsed_attributes)
end

When('I save all referentials') do
  RestClient.post url_for(path: "_referentials/save"), "", {:Authorization => "Token token=#{$adminToken}"}
end

When(/I reload the referential "([^"]+)"$/) do |referential|
  RestClient.post url_for(path: "_referentials/#{referential}/reload"), "", {content_type: :json, :Authorization => "Token token=#{$adminToken}"}
end
