require 'rest-client'
require 'json'

def url(referential = "test")
  referential ||= "test"
  "#{$server}/#{referential}/stop_areas"
end

def model_attributes table
  attributes = table.rows_hash
  if attributes["ObjectIds"]
    attributes["ObjectIds"] = JSON.parse("{#{attributes["ObjectIds"]}}")
  end
  attributes
end

Given(/^a StopArea exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  RestClient.post url(referential), model_attributes(stopArea).to_json, {content_type: :json}
end

When(/^a StopArea is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
  if referential.nil?
    step "a StopArea exists with the following attributes:", stopArea
  else
    step "a StopArea exists in Referential \"#{referential}\" with the following attributes:", stopArea
  end
end

When(/^the StopArea "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, objectid, referential|
  response = RestClient.get url(referential)
  responseArray = JSON.parse(response.body)
  expectedStopArea = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == kind && o["Value"] == objectid }}

  RestClient.delete "#{url(referential)}/#{expectedStopArea["Id"]}"
end

Then(/^one StopArea(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, stopArea|
  response = RestClient.get url(referential)
  responseArray = JSON.parse(response.body)

  stopAreaHash = model_attributes(stopArea)
  objectidkind = stopAreaHash["ObjectIds"].keys.first
  objectid_value = stopAreaHash["ObjectIds"][objectidkind]

  expectedName = responseArray.find{|a| a["Name"] == stopAreaHash["Name"]}
  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedName).not_to be_nil
  expect(expectedAttr).not_to be_nil
end


Then(/^a StopArea "([^"]+)":"([^"]+)" should (not )?exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get url(referential)
  responseArray = JSON.parse(response.body)
  expectedStopArea = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == kind && o["Value"] == objectid }}

  if condition.nil?
    expect(expectedStopArea).not_to be_nil
  else
    expect(expectedStopArea).to be_nil
  end
end
