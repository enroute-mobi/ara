def partners_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner'))
end

def getFirstPartner()
  response = RestClient.get partners_path(referential: "test"), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)
  response_array[0]["Id"]
end

Given(/^a Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\] and the following settings:$/) do |slug, referential, connectors, settings|
	attributes = {"slug" => slug, "connectorTypes" => connectors.split(',').map(&:strip), "settings" => settings.rows_hash}

  begin
	  RestClient.post partners_path(referential: referential), attributes.to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  rescue RestClient::ExceptionWithResponse => err
    puts err.response.body
    raise err
  end
end

Then(/^one Partner(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get partners_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)

  parsed_attributes = model_attributes(attributes)
  found_value = response_array.find{|a| a["Id"] == parsed_attributes["Id"]}

  expect(found_value).not_to be_nil

  expect(found_value).to include(parsed_attributes)
end


When(/^a Subscription exist (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, attributes|
  path = partners_path(referential: referential) + "/" + getFirstPartner() + "/subscriptions"
  response = RestClient.post path,  model_attributes(attributes).to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
end
