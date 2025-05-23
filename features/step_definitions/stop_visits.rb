def stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit'))
end

def scheduled_stop_visit_path(attributes = {})
  url_for_model(attributes.merge(resource: 'scheduled_stop_visit'))
end

def stop_visit_path(id, attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit', id: id))
end

Given(/^a ScheduledStopVisit exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_visit|
  response = RestClient.post scheduled_stop_visit_path(referential: referential), model_attributes(stop_visit).to_json, {content_type: :json, :Authorization => "Token token=#{$token}" }
  debug response.body
end

# Given(/^a StopVisit exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, stop_visit|
#   referential = find_referential(slug)
#   stop_visit = referential.stop_visits.create(model_attributes(stop_visit).transform_keys { |key| key.to_s.underscore.to_sym })

#   raise 'Cannot create stopArea' unless stop_visit.save
# end

# When(/^a StopVisit is created (?:in Referential "([^"]+)" )?with the following attributes:$/) do |referential, stop_visit|
#   if referential.nil?
#     step "a StopVisit exists with the following attributes:", stop_visit
#   else
#     step "a StopVisit exists in Referential \"#{referential}\" with the following attributes:", stop_visit
#   end
# end

# When(/^the StopVisit "([^"]*)"(?: in Referential "([^"]+)")? is edited with the following attributes:$/) do |code, slug, attributes|
#   stop_visit = find_model(slug, :stop_visits, code)
#   stop_visit.model_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

#   raise 'Cannot update stop_visit' unless stop_visit.save
# end

# Then(/^the StopVisit "([^"]*)" has the following attributes:$/) do |identifier, attributes|
#   stop_visit = find_model(nil, :stop_visits, identifier)
#   expect(stop_visit).not_to be_nil

#   matcher_attributes(attributes, stop_visit)
# end

# Then(/^one StopVisit has the following attributes:$/) do |attributes|
#   referential = find_referential(nil)

#   check_attributes(referential, :stop_visits, attributes)
# end

# Then(/^a StopVisit "([^"]+)":"([^"]+)" should( not)? exist(?: in Referential "([^"]+)")?$/) do |kind, value, condition, slug|
#   stop_visit = find_model(slug, :stop_visits, "#{kind}:#{value}")

#   if condition.nil?
#     expect(stop_visit).not_to be_nil
#   else
#     expect(stop_visit).to be_nil
#   end
# end
