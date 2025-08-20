def show_me(model_type, partner_name, slug = "test")
  models = referential_models(slug, model_type.to_sym)

  puts JSON.pretty_generate(models.map(&:api_attributes))
end

Given(/^show me ara subscriptions for partner "([^"]+)"?$/) do |partner|
  url = url_for_model({resource: 'partner', partner_name: "#{partner}", model: 'subscriptions'})
  response = RestClient.get(url, {content_type: :json, :Authorization => "Token token=#{$token}"})
  puts JSON.pretty_generate(JSON.parse(response))
end

Then(/^show me ara (vehicle_journeys|stop_areas|stop_area_groups|stop_visits|lines|line_groups|vehicles|partners|operators|scheduled_stop_visits|subscriptions|situations|facilities)$/) do |model_type|
  show_me(model_type, nil)
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
