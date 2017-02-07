require 'webrick'
require 'json'

def create_webrick_server(address, response)
	puts "start server"
	@webrick_server = WEBrick::HTTPServer.new(Port: 8090)
	@webrick_requests = []

	@webrick_server.mount_proc '/' do |req, res|
		body = req.body

		if req.body =~ /ns7:CheckStatus/
			res.body = %Q{
				<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:CheckStatusResponse xmlns:ns3="http://www.siri.org.uk/siri"
                             xmlns:ns4="http://www.ifopt.org.uk/acsb"
                             xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                             xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                             xmlns:ns7="http://scma/siri"
                             xmlns:ns8="http://wsdl.siri.org.uk"
                             xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <CheckStatusAnswerInfo>
        <ns3:ResponseTimestamp>2016-09-22T07:58:34.000+02:00</ns3:ResponseTimestamp>
        <ns3:ProducerRef>NINOXE:default</ns3:ProducerRef>
        <ns3:Address>http://appli.chouette.mobi/siri_france/siri</ns3:Address>
        <ns3:ResponseMessageIdentifier>c464f588-5128-46c8-ac3f-8b8a465692ab</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>CheckStatus:Test:0</ns3:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <ns3:Status>true</ns3:Status>
        <ns3:ServiceStartedTime>2016-09-22T03:30:32.000+02:00</ns3:ServiceStartedTime>
      </Answer>
      <AnswerExtension />
    </ns8:CheckStatusResponse>
  </S:Body>
</S:Envelope>
			}
		else
			@webrick_requests << req
		  res.body = response
		end

		res.content_type = "text/xml"
	end
	
	Thread.start do
		@webrick_server.start
	end
end

Given(/^a SIRI server waits GetStopMonitoring request on "([^"]*)" to respond with$/) do |address, response|
	create_webrick_server(address, response)
end

def partners_path(attributes = {})
  url_for_model(attributes.merge(resource: 'partner'))
end

Given(/^a Partner "([^"]*)" exists with connectors \[([^"\]]*)\] and the following settings:$/) do |partner, connectors, settings|
	attributes = {"slug" => partner, "connectorTypes" => connectors.split(',').map(&:strip), "settings" => model_attributes(settings)}
	puts partners_path, attributes.inspect
	puts RestClient.post partners_path, attributes.to_json, {content_type: :json, accept: :json}
end

def time_path(action = "")
  base_url = url_for(path: "_time")
  base_url += "/#{action}" unless action.empty?
  base_url
end

When(/^a minute has passed$/) do
	puts time_path("advance"), { "duration" => "60s" }.to_json
	RestClient.post(time_path("advance"), { "duration" => "60s" }.to_json)
end

When(/^the SIRI server should have receive a GetStopMonitoring request$/) do
	try_count = 0
	while @webrick_requests.empty?
		try_count += 1
		raise "No received request" if try_count > 10

		sleep 0.5
	end
end

def stop_visits_path(attributes = {})
  url_for_model(attributes.merge(resource: 'stop_visit'))
end

Then(/^one StopVisit has the following attributes:$/) do |attributes|
	response = RestClient.get stop_visits_path
	responseArray = JSON.parse(response.body)

	attributes = model_attributes(attributes)

	objectidkind = attributes["ObjectIds"].keys.first
  objectid_value = attributes["ObjectIds"][objectidkind]

	expectedName = responseArray.find{|a| a["Name"] == attributes["Name"]}
  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedName).not_to be_nil
	expect(expectedAttr).not_to be_nil
end

def vehicle_journey_path(attributes = {})
	url_for_model(attributes.merge(resource: 'vehicle_journey'))
end

Then(/^one VehicleJourney has the following attributes:$/) do |attributes|
	response = RestClient.get vehicle_journey_path
	responseArray = JSON.parse(response.body)

	attributes = model_attributes(attributes)

	objectidkind = attributes["ObjectIds"].keys.first
  objectid_value = attributes["ObjectIds"][objectidkind]

  expectedAttr = responseArray.find{|a| a["ObjectIDs"].find{|o| o["Kind"] == objectidkind && o["Value"] == objectid_value }}

  expect(expectedAttr).not_to be_nil
end