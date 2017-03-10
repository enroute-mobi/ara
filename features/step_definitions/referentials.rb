def referentials_path
  url_for(path: "_referentials")
end

def referential_path(id)
  url_for(path: "_referentials/#{id}")
end

Given(/^a Referential "([^"]+)" exists$/) do |referential|
  attributes = {
    slug: referential
  }
  RestClient.post referentials_path, attributes.to_json, {content_type: :json}
end

Given(/^a Referential "([^"]*)" exists with the following settings:$/) do |referential, settings|
  attributes = {
    slug: referential
  }
  attributes[:settings] = settings.rows_hash if settings
  RestClient.post referentials_path, attributes.to_json, {content_type: :json}
end

When(/^a Referential "([^"]+)" is created$/) do |referential|
  step "a Referential \"#{referential}\" exists"
end

When(/^the Referential "([^"]+)" is destroyed$/) do |referential|
  response = RestClient.get referentials_path
  responseHash = JSON.parse(response.body)

  id = responseHash.find{|a| a["Slug"] == referential}["Id"]
  RestClient.delete referential_path(id)
end

Then(/^a Referential "([^"]+)" should (not )?exist$/) do |referential, condition|
  response = RestClient.get referentials_path
  responseHash = JSON.parse(response.body)

  if condition.nil?
    expect(responseHash.find{|a| a["Slug"] == referential}).not_to be_nil
  else
    expect(responseHash.find{|a| a["Slug"] == referential}).to be_nil
  end
end
