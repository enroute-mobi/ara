def operators_path(attributes = {})
  url_for_model(attributes.merge(resource: 'operator'))
end

def operator_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'operator', id: id))
end


Given(/^an Operator exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, operator|
  RestClient.post operators_path(referential: referential), model_attributes(operator).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
end

# When(/^a Operator is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stopArea|
#   if referential.nil?
#     step "a Operator exists with the following attributes:", stopArea
#   else
#     step "a Operator exists in Referential \"#{referential}\" with the following attributes:", stopArea
#   end
# end

# When(/^the Operator "([^"]*)" is edited with the following attributes:$/) do |identifier, attributes|
#   RestClient.put operator_path(identifier), model_attributes(attributes).to_json, {content_type: :json, :Authorization => "Token token=#{$token}"}
# end

# Then(/^the Operator "([^"]*)" has the following attributes:$/) do |identifier, attributes|
#   response = RestClient.get operator_path(identifier), {content_type: :json, :Authorization => "Token token=#{$token}"}
#   stopVisitAttributes = api_attributes(response.body)
#   expect(stopVisitAttributes).to include(model_attributes(attributes))
# end

# Then(/^one Operator has the following attributes:$/) do |attributes|
#   response = RestClient.get operators_path, {content_type: :json, :Authorization => "Token token=#{$token}"}
#   response_array = JSON.parse(response.body)

#   called_method = has_attributes(response_array, attributes)

#   expect(called_method).to be_truthy
# end


# Then(/^a Operator "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, referential|
#   response = RestClient.get(operator_path("#{kind}:#{value}", referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}){|response, request, result| response }

#   if condition.nil?
#     expect(response.code).to eq(200)
#   else
#     expect(response.code).to eq(404)
#     expect(response.body).to include("Stop visit not found: #{kind}:#{value}")
#   end
# end
