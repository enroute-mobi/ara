def partner_templates_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner_template'))
end

def getFirstPartnerTemplate()
  response = RestClient.get partner_templates_path(referential: "test"), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)
  response_array[0]["Id"]
end

Given(/^a Partner template "([^"]*)" exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, referential, attributes|
  attrs = model_attributes(attributes)
  attrs["slug"] = slug

  begin
    RestClient.post partner_templates_path(referential: referential), attrs.to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  rescue RestClient::ExceptionWithResponse => err
    puts err.response.body
    raise err
  end
end

When(/^the Partner template "([^"]*)"(?: in Referential "([^"]+)")? is updated with the following attributes:$/) do |slug, referential, attributes|
  attrs = attributes.rows_hash

  path = partner_templates_path + '/' + slug
  begin
      RestClient.put path, attrs.to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  rescue RestClient::ExceptionWithResponse => err
    puts err.response.body
    raise err
  end
end

Then(/^one Partner template(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get partner_templates_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)

  parsed_attributes = model_attributes(attributes)
  found_value = response_array.find{|a| a["Slug"] == parsed_attributes["Slug"]}

  expect(found_value).not_to be_nil

  expect(found_value).to include(parsed_attributes)
end
