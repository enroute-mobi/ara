Given(/^a StopArea Group exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, stop_area_group|
  referential = find_referential(slug)
  stop_area_group = referential.stop_area_groups.create(model_attributes(stop_area_group).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create stop_area group' unless stop_area_group.save
end

When(/^a StopArea Group is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_area_group|
  if referential.nil?
    step "a StopArea Group exists with the following attributes:", stop_area_group
  else
    step "a StopArea Group exists in Referential \"#{referential}\" with the following attributes:", stop_area_group
  end
end

Then(/^the StopArea Group "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  stop_area_group = find_model(nil, :stop_area_groups, identifier)
  expect(stop_area_group).not_to be_nil

  matcher_attributes(attributes, stop_area_group)
end
