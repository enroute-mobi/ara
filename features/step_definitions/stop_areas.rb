Given(/^a StopArea exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, stop_area|
  referential = find_referential(slug)
  stop_area = referential.stop_areas.create(model_attributes(stop_area).transform_keys { |key| key.to_s.underscore.to_sym })

  raise 'Cannot create stopArea' unless stop_area.save
end

When(/^a StopArea is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_area|
  if referential.nil?
    step "a StopArea exists with the following attributes:", stop_area
  else
    step "a StopArea exists in Referential \"#{referential}\" with the following attributes:", stop_area
  end
end

When(/^the StopArea "([^"]+)":"([^"]+)"(?: in Referential "([^"]+)")? is destroyed$/) do |kind, value, slug|
  stop_area = find_model(slug, :stop_areas, "#{kind}:#{value}")

  raise 'Cannot destroy stopArea' unless stop_area.destroy
end

Then(/^one StopArea(?: in Referential "([^"]+)")? has the following attributes:$/) do |slug, attributes|
  referential = find_referential(slug)

  check_attributes(referential, :stop_areas, attributes)
end

Then(/^a StopArea "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, slug|
  stop_area = find_model(slug, :stop_areas, "#{kind}:#{value}")

  if condition.nil?
    expect(stop_area).not_to be_nil
  else
    expect(stop_area).to be_nil
  end
end
