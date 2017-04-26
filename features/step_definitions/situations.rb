def situations_path(attributes = {})
  url_for_model(attributes.merge(resource: 'situation'))
end

def situation_path(id, attributes = {})
  path = url_for_model(attributes.merge(resource: 'situation', id: id))
end


Then(/^one Situation(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get situations_path(referential: referential)
  response_array = JSON.parse(response.body)

  called_method = has_attributes(response_array, attributes)

  expect(called_method).to be_truthy
end