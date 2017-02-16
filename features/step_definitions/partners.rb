def partners_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner'))
end

Given(/^a Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\] and the following settings:$/) do |slug, referential, connectors, settings|
	attributes = {"slug" => slug, "connectorTypes" => connectors.split(',').map(&:strip), "settings" => settings.rows_hash}
	RestClient.post partners_path(referential: referential), attributes.to_json, {content_type: :json, accept: :json}
end

# Given(/^a local Partner "([^"]*)" exists (?:in Referential "([^"]+)" )?with connectors \[([^"\]]*)\]$/) do |slug, referential, connectors, settings|
#   # FIXME after #2560 and #2561
#   connector_aliases = {
#     "siri-check-status-server" => "siri-check-status-client",
#     "siri-stop-monitoring-request-broadcaster" => "siri-stop-monitoring-request-collector"
#   }

#   connectors = connectors.split(',').map(&:strip).map do |connector|
#     connectors = connector_aliases.fetch(connector, connector)
#   end

#   settings = settings.rows_hash

#   step %Q{a Partner "#{slug}" exists with connectors [#{connectors.join(',')}] and the following settings:}, table(%Q{
#       | local_credential     | #{settings["local_credential"]} |
#       | remote_url           | http://localhost:8090           |
#       | remote_credential    | test                            |
#       | remote_objectid_kind | internal                        |
#   })
# end
