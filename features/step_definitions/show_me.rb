def show_me(model_type, referential = "test")
  response = RestClient.get send("#{model_type}_path", referential: referential), {content_type: :json, :Authorization => "Token token=#{$token}"}
  puts JSON.pretty_generate(JSON.parse(response))
end

Given(/^I see edwig (vehicle_journeys|stop_areas|stop_visits|lines|partners|operators)$/) do |model_type|
  show_me model_type
end

Then(/^show me edwig (vehicle_journeys|stop_areas|stop_visits|lines|partners|operators)$/) do |model_type|
  show_me model_type
end

def show_me_time
  time = Time.parse(JSON.parse(RestClient.get(time_path).body)["time"])
  puts "Edwig time is #{time}"
end

Given(/^I see edwig time$/) do
  show_me_time
end

Then(/^show me edwig time$/) do
  show_me_time
end
