def set_default_generators!(settings)
  default_generators_map = {
    "generators.message_identifier" => "RATPDev:Message::%{uuid}:LOC",
    "generators.response_message_identifier" => "RATPDev:ResponseMessage::%{uuid}:LOC",
    "generators.data_frame_identifier" => "RATPDev:DataFrame::%{id}:LOC",
    "generators.reference_identifier" => "RATPDev:%{type}::%{default}:LOC",
    "generators.reference_stop_area_identifier" => "RATPDev:StopPoint:Q:%{default}:LOC"
  }

  default_generators_map.each do |id, format_string|
    settings[id] ||= format_string
  end
end

Given(/^a (SIRI )?Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\] and the following settings:$/) do |siri, slug, referential_slug, connectors, settings_table|
  settings = settings_table.rows_hash
  set_default_generators!(settings) if siri

  referential = find_referential(referential_slug)
  partner = referential.partners.create(
    slug: slug,
    connector_types: connectors.split(',').map(&:strip),
    settings: settings
  )

  raise "Cannot create Partner: #{partner.errors}" unless partner.save
end

Given('the Partner {string} is updated with the following settings:') do |slug, settings_table|
  partner = find_referential('test').partners.find(slug)
  partner.settings = settings_table.rows_hash

  raise "Cannot update Partner: #{partner.errors}" unless partner.save
end

Then(/^one Partner(?: in Referential "([^"]+)")? has the following attributes:$/) do |referential_slug, attributes|
  partners = find_referential(referential_slug).partners.all
  parsed_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  found = partners.find { |p| p.id == parsed_attributes[:id] }
  expect(found).not_to be_nil
  expect(found).to have_attributes(parsed_attributes)
end

Then(/^the Partner "([^"]+)" in the Referential "([^"]+)" has the operational status (up|down|unknown)/) do |slug, referential_slug, status|
  partner = find_referential(referential_slug).partners.find(slug)
  expect(partner.status).to eq(status)
end

When(/^a Subscription exist (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential_slug, attributes|
  partner = find_referential(referential_slug).partners.all.first
  sub_attrs = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }
  sub = partner.subscriptions.create(sub_attrs)
  sub.save
end

Then(/^No Subscription exists with the following attributes:$/) do |attributes|
  partner = find_referential('test').partners.all.first
  subs = partner.subscriptions.all
  parsed_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  expect(subs).not_to include(an_object_having_attributes(parsed_attributes))
end

Then(/^one Subscription exists with the following attributes:$/) do |attributes|
  partner = find_referential('test').partners.all.first
  subs = partner.subscriptions.all

  attrs = attributes.rows_hash

  if (subscribed_at = attrs.delete("Resources[0]/SubscribedAt"))
    subscribed_at = (a_value > $1) if %r{^> (.*)$} =~ subscribed_at
    attrs["Resources"] = a_collection_including(a_hash_including(subscribed_at: subscribed_at))
  end

  parsed_attributes = attrs.transform_keys { |key| key.to_s.underscore.to_sym }

  expect(subs).to include(an_object_having_attributes(parsed_attributes))
end

Then(/^Subscriptions exist with the following resources:$/) do |attributes|
  partner = find_referential('test').partners.all.first
  subs = partner.subscriptions.all

  subscription_codes = subs.flat_map(&:resources)
                           .map { |r| r[:reference]["Code"] }

  attributes.to_hash.map { |v| { v[0] => v[1] } }.each do |expected_subscription|
    expect(subscription_codes).to include(expected_subscription)
  end
end

Then(/^No Subscriptions exist with the following resources:$/) do |attributes|
  partner = find_referential('test').partners.all.first
  subs = partner.subscriptions.all

  first_sub_resources = (subs.first&.resources || []).map { |r| r[:reference]["Code"] }

  attributes.to_hash.map { |v| { v[0] => v[1] } }.each do |expected_subscription|
    expect(first_sub_resources).not_to include(expected_subscription)
  end
end

Then(/^no Subscription exists/) do
  partner = find_referential('test').partners.all.first
  subs = partner.subscriptions.all

  expect(subs).to be_empty
end

When(/^I wait that a Subscription has been created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential_slug, attributes|
  partner = find_referential(referential_slug).partners.all.first

  subs = []
  while subs.empty?
    step "10 seconds have passed"

    subs = partner.subscriptions.all.reject do |sub|
      (sub.resources || []).all? { |r| r[:subscribed_at] == "0001-01-01T00:00:00Z" }
    end
  end

  parsed_attributes = attributes.rows_hash.transform_keys { |key| key.to_s.underscore.to_sym }
  expect(subs).to include(an_object_having_attributes(parsed_attributes))
end
