Given(/^a Line exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, line|
  referential = find_referential(slug)
  line = referential.lines.create(model_attributes(line).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create line' unless line.save
end

When(/^a Line is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, line|
  if referential.nil?
    step 'a Line exists with the following attributes:', line
  else
    step "a Line exists in Referential \"#{referential}\" with the following attributes:", line
  end
end

When(/^the Line "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, value, slug|
  line = find_model(slug, :lines, "#{kind}:#{value}")

  raise 'Cannot destroy line' unless line.destroy
end

Then(/^one Line(?: in Referential "([^"]+)")? has the following attributes:$/) do |slug, attributes|
  referential = find_referential(slug)

  check_attributes(referential, :lines, attributes)
end

Then(/^a Line "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, slug|
  line = find_model(slug, :lines, "#{kind}:#{value}")

  if condition.nil?
    expect(line).not_to be_nil
  else
    expect(line).to be_nil
  end
end
