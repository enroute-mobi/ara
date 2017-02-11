def partners_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner'))
end

Given(/^a Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\] and the following settings:$/) do |slug, referential, connectors, settings|
	attributes = {"slug" => slug, "connectorTypes" => connectors.split(',').map(&:strip), "settings" => settings.rows_hash}
	RestClient.post partners_path(referential: referential), attributes.to_json, {content_type: :json, accept: :json}
end
