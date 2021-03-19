When('I send a GTFS-RT request to the Referential {string} without token') do |string|
  @gtfs_response = GTFS::Realtime.get gtfs_url(referential: string)
end

When('I send a GTFS-RT request to the Referential {string} with token {string}') do |referential_slug, partner_token|
  headers = { "Authorization" => "Token token=#{partner_token}" }
  @gtfs_response = GTFS::Realtime.get gtfs_url(referential: referential_slug), headers
end

Then('I should receive a GTFS-RT response') do
  expect(@gtfs_response).to be_an_instance_of(GTFS::Realtime::FeedMessage)
end

Then('I should not receive a GTFS-RT but an unauthorized client error status') do
  expect(@gtfs_response).to be_an_instance_of(Net::HTTPUnauthorized)
end

Then('this GTFS-RT response should contain a Vehicle Position with these attributes:') do |vehicle_position_attributes|
  vehicle_position_attributes = vehicle_position_attributes.rows_hash
  debug @gtfs_response.vehicle_positions.inspect
  expect(@gtfs_response.vehicle_positions).to include(an_object_having_attributes(vehicle_position_attributes))
end

Then('this GTFS-RT response should not contain Vehicle Positions') do
  debug @gtfs_response.vehicle_positions.inspect
  expect(@gtfs_response.vehicle_positions).to be_empty
end
