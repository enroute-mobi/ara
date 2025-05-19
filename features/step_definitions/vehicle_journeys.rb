Given(/^a VehicleJourney exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, attributes|
  referential = find_referential(slug)
  vehicle_journey = referential.vehicle_journeys.create(model_attributes(attributes).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create VehicleJourney' unless vehicle_journey.save
end

When(/^the VehicleJourney "([^"]*)"(?: in Referential "([^"]+)")? is edited with the following attributes:$/) do |identifier, slug, attributes|
  vehicle_journey = find_model(slug, :vehicle_journeys, identifier)
  vehicle_journey.model_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  raise 'Cannot update VehicleJourney' unless vehicle_journey.save
end

Then(/^the VehicleJourney "([^"]*)" has the following attributes:$/) do |identifier, attributes|
  vehicle_journey = find_model(nil, :vehicle_journeys, identifier)
  expect(vehicle_journey).not_to be_nil

  matcher_attributes(attributes, vehicle_journey)
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
  referential = find_referential(nil)

  check_attributes(referential, :vehicle_journeys, attributes)
end

Then(/^a VehicleJourney "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, slug|
  vehicle_journey = find_model(slug, :vehicle_journeys, "#{kind}:#{value}")

  if condition.nil?
    expect(vehicle_journey).not_to be_nil
  else
    expect(vehicle_journey).to be_nil
  end
end
