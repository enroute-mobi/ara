Given(/^a Vehicle exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, attributes|
  referential = find_referential(slug)
  vehicle = referential.vehicles.create(model_attributes(attributes).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create vehicle' unless vehicle.save
end

When(/^a Vehicle is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, vehicle|
  if referential.nil?
    step "a Vehicle exists with the following attributes:", vehicle
  else
    step "a Vehicle exists in Referential \"#{referential}\" with the following attributes:", vehicle
  end
end

When(/^the Vehicle "([^"]*)"(?: in Referential "([^"]+)")? is edited with the following attributes:$/) do |code, slug, attributes|
  vehicle = find_model(slug, :vehicles, code)
  vehicle.model_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  raise 'Cannot update vehicle' unless vehicle.save
end

Then(/^one Vehicle has the following attributes:$/) do |attributes|
  referential = find_referential(nil)

  check_attributes(referential, :vehicles, attributes)
end

Then(/^a Vehicle "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, slug|
  vehicle = find_model(slug, :vehicle, "#{kind}:#{value}")

  if condition.nil?
    expect(vehicle).not_to be_nil
  else
    expect(vehicle).to be_nil
  end
end
