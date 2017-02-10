def time_path(action = "")
  base_url = url_for(path: "_time")
  base_url += "/#{action}" unless action.empty?
  base_url
end

When(/^a minute has passed$/) do
	RestClient.post(time_path("advance"), { "duration" => "60s" }.to_json)
	sleep 1 # don't blame me
end
