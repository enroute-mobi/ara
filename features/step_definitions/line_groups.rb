Given(/^a Line Group exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, line_group|
  referential = find_referential(slug)
  line_group = referential.line_groups.create(model_attributes(line_group).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create line group' unless line_group.save
end

When(/^a Line Group is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line_group|
  if referential.nil?
    step "a Line Group exists with the following attributes:", line_group
  else
    step "a Line Group exists in Referential \"#{referential}\" with the following attributes:", line_group
  end
end

Then(/^the Line Group "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  line_group = find_model(nil, :line_groups, identifier)
  expect(line_group).not_to be_nil

  matcher_attributes(attributes, line_group)
end
