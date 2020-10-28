def time_path(action = "")
  base_url = url_for(path: "_time")
  base_url += "/#{action}" unless action.empty?
  base_url
end

When(/^a minute has passed$/) do
	RestClient.post(time_path("advance"), { "duration" => "60s" }.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"})
	sleep 1 # don't blame me
end

When(/^(\d+) minutes have passed$/) do |count|
	RestClient.post(time_path("advance"), { "duration" => "#{count.to_i * 60}s" }.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"})
	sleep 1 # don't blame me
end

When(/^(\d+) seconds have passed$/) do |count|
	RestClient.post(time_path("advance"), { "duration" => "#{count.to_i}s" }.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"})
	sleep 1 # don't blame me
end

When(/^the time is "([^"]*)"$/) do |expected_time|
	getTime = RestClient.get(time_path, {content_type: :json, :Authorization => "Token token=#{$adminToken}"}).body
	splitTime = getTime.split(' ')

	araTime = Time.parse(splitTime[2])
	expectedTime = Time.parse(expected_time)

	duration =  expectedTime - araTime

	RestClient.post(time_path("advance"), { "duration" => "#{duration}s" }.to_json, {content_type: :json, :Authorization => "Token token=#{$adminToken}"})
	sleep 1 # don't blame me
end
