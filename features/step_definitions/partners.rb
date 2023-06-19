def partners_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner'))
end

def getFirstPartner()
  response = RestClient.get partners_path(referential: "test"), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)
  response_array[0]["Id"]
end

def set_default_generators! attributes
  default_generators_map = {
    "generators.message_identifier" => "RATPDev:Message::%{uuid}:LOC",
    "generators.response_message_identifier" => "RATPDev:ResponseMessage::%{uuid}:LOC",
    "generators.data_frame_identifier" => "RATPDev:DataFrame::%{id}:LOC",
    "generators.reference_identifier" => "RATPDev:%{type}::%{default}:LOC",
    "generators.reference_stop_area_identifier" => "RATPDev:StopPoint:Q:%{default}:LOC"
  }

  default_generators_map.each do |id, formatString|
    attributes["settings"][id] = formatString if attributes["settings"][id].nil?
  end
end

Given(/^a (SIRI )?Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\] and the following settings:$/) do |siri, slug, referential, connectors, settings|
	attributes = {"slug" => slug, "connectorTypes" => connectors.split(',').map(&:strip), "settings" => settings.rows_hash}
  # Set default generators to avoid updating all cucumber tests
  if siri
    set_default_generators!(attributes)
  end

  begin
	  RestClient.post partners_path(referential: referential), attributes.to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  rescue RestClient::ExceptionWithResponse => err
    puts err.response.body
    raise err
  end
end

Given('the Partner {string} is updated with the following settings:') do |slug, settings|
	attributes = {"slug" => slug, "settings" => settings.rows_hash}
  path = partners_path + '/' + slug
	RestClient.put path, attributes.to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
end

Then(/^one Partner(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential, attributes|
  response = RestClient.get partners_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)

  parsed_attributes = model_attributes(attributes)
  found_value = response_array.find{|a| a["Id"] == parsed_attributes["Id"]}

  expect(found_value).not_to be_nil

  expect(found_value).to include(parsed_attributes)
end

Then(/^the Partner "([^"]+)" in the Referential "([^"]+)" has the operational status (up|down|unknown)/) do |slug, referential, status|
  response = RestClient.get partners_path(referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  response_array = api_attributes(response.body)
  partner = response_array.find { |partner| partner['Slug'] == slug }
  operational_status = partner['PartnerStatus']['OperationnalStatus']

  expect(operational_status).to eq(status)
end

When(/^a Subscription exist (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, attributes|
  path = partners_path(referential: referential) + "/" + getFirstPartner() + "/subscriptions"
  response = RestClient.post path,  model_attributes(attributes).to_json, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}

  debug response.body
end

Then(/^No Subscription exists with the following attributes:$/) do |attributes|
  path = partners_path + '/' + getFirstPartner() + '/subscriptions'
  response = RestClient.get path, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  attributes = attributes.rows_hash

  expect(response_array).not_to include(a_hash_including(attributes))
end

Then(/^one Subscription exists with the following attributes:$/) do |attributes|
  path = partners_path + '/' + getFirstPartner() + '/subscriptions'
  response = RestClient.get path, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  attributes = attributes.rows_hash
  # TODO: build complex matcher from attributes ...
  if (subscribed_at = attributes.delete("Resources[0]/SubscribedAt"))
    if %r{^> (.*)$} =~ subscribed_at
      subscribed_at = (a_value > $1)
    end
    attributes["Resources"] = a_collection_including(a_hash_including("SubscribedAt" => subscribed_at))
  end

  expect(response_array).to include(a_hash_including(attributes))
end

Then(/^Subscriptions exist with the following resources:$/) do |attributes|
  path = partners_path + '/' + getFirstPartner + '/subscriptions'
  response = RestClient.get path, { content_type: :json, accept: :json, :Authorization => "Token token=#{$token}" }
  response_array = JSON.parse(response.body)

  subscriptions = response_array.map { |s| s['Resources'] }
                    .flatten
                    .map { |r| r['Reference']['ObjectId'] }

  attributes.to_hash.map { |v| { v[0] => v[1] } }.each do |expected_subscription|
    expect(subscriptions).to include(expected_subscription)
  end
end

Then(/^No Subscriptions exist with the following resources:$/) do |attributes|
  path = partners_path + '/' + getFirstPartner + '/subscriptions'
  response = RestClient.get path, { content_type: :json, accept: :json, :Authorization => "Token token=#{$token}" }
  response_array = JSON.parse(response.body)

  subscriptions = response_array.first['Resources'].map { |r| r['Reference']['ObjectId'] }

  attributes.to_hash.map { |v| { v[0] => v[1] } }.each do |expected_subscription|
    expect(subscriptions).not_to include(expected_subscription)
  end
end

Then(/^no Subscription exists/) do
  path = partners_path + '/' + getFirstPartner() + '/subscriptions'
  response = RestClient.get path, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
  response_array = JSON.parse(response.body)

  expect(response_array).to eq([])
end

When(/^I wait that a Subscription has been created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, attributes|
  path = partners_path(referential: referential) + "/" + getFirstPartner() + "/subscriptions"

  retry_count = 0
  response_array = []
  while response_array.empty?
    step "10 seconds have passed"

    response  = RestClient.get path, {content_type: :json, accept: :json, :Authorization => "Token token=#{$token}"}
    response_array = api_attributes(response.body)

    # We're ignoring pending subscriptions
    response_array.delete_if do |subscription|
      subscription["Resources"].delete_if do |resource|
        resource["SubscribedAt"] == "0001-01-01T00:00:00Z"
      end
      subscription["Resources"].empty?
    end

    retry_count += 1
  end

  expect(response_array).to include(a_hash_including(attributes.rows_hash))
end
