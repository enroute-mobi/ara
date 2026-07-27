def show_me(model_type, slug = "test")
  models = referential_models(slug, model_type.to_sym)

  puts JSON.pretty_generate(models.map(&:api_attributes))
end

Given(/^show me ara subscriptions for partner "([^"]+)"?$/) do |partner|
  partner_obj = find_referential('test').partners.find(partner)
  subs = partner_obj.subscriptions.all
  puts JSON.pretty_generate(subs.map { |s|
    {
      subscription_ref: s.subscription_ref,
      external_id: s.external_id,
      kind: s.kind,
      resources: s.resources,
      subscription_options: s.subscription_options
    }.compact
  })
end

Then(/^show me ara (vehicle_journeys|stop_areas|stop_area_groups|stop_visits|lines|line_groups|vehicles|partners|operators|scheduled_stop_visits|subscriptions|situations|facilities)$/) do |model_type|
  show_me(model_type)
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
