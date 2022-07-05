def scheduled_stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'scheduled_stop_visit'))
end
