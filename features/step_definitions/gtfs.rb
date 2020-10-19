When('I send this GTFS-RT request to the Referential {string} without token') do |string|
  @gtfs_response = GTFS::Realtime.get gtfs_url(referential: string)
end

When('I send this GTFS-RT request to the Referential {string} with token {string}') do |referential_slug, partner_token|
  headers = { "Authorization" => "Token token=#{partner_token}" }
  @gtfs_response = GTFS::Realtime.get gtfs_url(referential: referential_slug), headers
end

Then('I should receive a GTFS-RT response') do
  expect(@gtfs_response).to be_an_instance_of(GTFS::Realtime::FeedMessage)
end

Then('I should not receive a GTFS-RT but an unauthorized client error status') do
  expect(@gtfs_response).to be_an_instance_of(Net::HTTPUnauthorized)
end
