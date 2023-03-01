def subscriptions_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner', model: 'subscriptions'))
end
