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
  expectedLine = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == kind && o["Value"] == objectid }}

  RestClient.delete line_path expectedLine["Id"]
end

Then(/^one Line(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, line|
  response = RestClient.get lines_path(referential: referential)
  responseArray = JSON.parse(response.body)

  lineHash = model_attributes(line)
  objectidkind = lineHash["ObjectIDs"].keys.first
  objectid_value = lineHash["ObjectIDs"][objectidkind]

  expectedName = responseArray.find{|a| a["Name"] == lineHash["Name"]}
  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedName).not_to be_nil
  expect(expectedAttr).not_to be_nil
end


Then(/^a Line "([^"]+)":"([^"]+)" should (not )?exist(?: in Referential "([^"]+)")?$/) do |kind, objectid, condition, referential|
  response = RestClient.get lines_path(referential: referential)
  responseArray = JSON.parse(response.body)
  expectedLine = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == kind && o["Value"] == objectid }}

  if condition.nil?
    expect(expectedLine).not_to be_nil
  else
    expect(expectedLine).to be_nil
  end
end