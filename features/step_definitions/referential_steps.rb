require 'rest-client'
require 'json'

url = "#{$server}/_referentials"

Given(/^a Referential "([^"]+)" exists$/) do |referential|
  RestClient.post url, "{\"slug\":\"#{referential}\"}", {content_type: :json}
end

When(/^a Referential "([^"]+)" is created$/) do |referential|
  step "a Referential \"#{referential}\" exists"
end

When(/^the Referential "([^"]+)" is destroyed$/) do |referential|
  response = RestClient.get url
  responseHash = JSON.parse(response.body)

  id = responseHash.find{|a| a["Slug"] == referential}["Id"]
  RestClient.delete "#{url}/#{id}"
end

Then(/^a Referential "([^"]+)" should (not )?exist$/) do |referential, condition|
  response = RestClient.get url
  responseHash = JSON.parse(response.body)

  if condition.nil?
    expect(responseHash.find{|a| a["Slug"] == referential}).not_to be_nil
  else
    expect(responseHash.find{|a| a["Slug"] == referential}).to be_nil
  end
end
