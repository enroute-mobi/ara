def stop_area_groups_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_area_group'))
end

def stop_area_group_path(id, attributes = {})
  ufrl_for_model(attributes.merge(resource: 'stop_area_groups', id: id))
end

Given(/^a StopArea Group exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopAreaGroup|
  response = RestClient.post stop_area_groups_path(referential: referential), model_attributes(stopAreaGroup).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
  debug response.body
end


When(/^a StopArea Group is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopAreaGroup|
  if referential.nil?
    step "a StopArea Group exists with the following attributes:", stopAreaGroup
  else
    step "a StopArea Group exists in Referential \"#{referential}\" with the following attributes:", stopAreaGroup
  end
end

Then(/^one StopArea Group(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get stop_area_groups_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)

  expect(response_array).to include(model_attributes(attributes))
end
