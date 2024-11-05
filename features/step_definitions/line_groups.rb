def line_groups_path(attributes = {})
  url_for_model(attributes.merge(resource: 'line_group'))
end

def line_group_path(id, attributes = {})
  ufrl_for_model(attributes.merge(resource: 'line_groups', id: id))
end

Given(/^a Line Group exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, lineGroup|
  response = RestClient.post line_groups_path(referential: referential), model_attributes(lineGroup).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
  debug response.body
end


When(/^a Line Group is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, lineGroup|
  if referential.nil?
    step "a Line Group exists with the following attributes:", lineGroup
  else
    step "a Line Group exists in Referential \"#{referential}\" with the following attributes:", lineGroup
  end
end

Then(/^one Line Group(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get line_groups_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)

  expect(response_array).to include(model_attributes(attributes))
end
