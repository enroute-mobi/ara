def show_me(model_type, partner_name, referential = "test")
  response = RestClient.get send("#{model_type}_path", referential: referential, partner_name: partner_name), {content_type: :json, :Authorization => "Token token=#{$token}"}
  puts JSON.pretty_generate(JSON.parse(response))
end

Given(/^I see ara (vehicle_journeys|stop_areas|stop_visits|lines|vehicles|partners|operators|scheduled_stop_visits)$/) do |model_type|
  show_me(model_type, nil)
end

Then(/^show me ara (vehicle_journeys|stop_areas|stop_area_groups|stop_visits|lines|line_groups|vehicles|partners|operators|scheduled_stop_visits|subscriptions|situations)(?: for partner "([^"]+)")?$/) do |model_type, partner_name|
  show_me(model_type, partner_name)
end

def show_me_time
  time = Time.parse(JSON.parse(RestClient.get(time_path).body)["time"])
  puts "Ara time is #{time}"
end

Given(/^I see ara time$/) do
  show_me_time
end

Then(/^show me ara time$/) do
  show_me_time
end
