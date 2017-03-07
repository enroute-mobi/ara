def show_me(model_type, referential = "test")
  response = RestClient.get send("#{model_type}_path", referential: referential)
  puts JSON.pretty_generate(JSON.parse(response))
end

Given(/^I see edwig (vehicle_journeys|stop_areas|stop_visits|lines|partners)$/) do |model_type|
  show_me model_type
end

Then(/^show me edwig (vehicle_journeys|stop_areas|stop_visits|lines|partners)$/) do |model_type|
  show_me model_type
end
