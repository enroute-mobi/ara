def siri_path(attributes = {})
  attributes = {
    referential: 'test'
  }.merge(attributes.delete_if { |k,v| v.nil? })

  url_for(attributes.merge(path: "siri"))
end

def send_siri_request(request, attributes = {})
  response = RestClient.post siri_path(attributes), request, {content_type: :xml}
  save_siri_messages request: request, response: response.body
  @last_siri_request = request
  @last_siri_response = response.body
end

def save_siri_messages(messages = {})
  return unless ENV['SIRI_DEBUG']

  @siri_timestamp ||= Time.now.strftime("%Y%m%d%H%M%S")
  @siri_message_id ||= 0
  @siri_message_id += 1

  messages.each do |type, content|
    file = "log/siri-message-#{@siri_timestamp}-#{@siri_message_id}-#{type}"
    File.write file, content, mode: "wb"
  end
end

Given(/^a SIRI server (?:"([^"]*)" )?waits (GetStopMonitoring) request on "([^"]*)" to respond with$/) do |name, message_type, url, response|
  name ||= "default"
  SIRIServer.create(name, url).expect_request(message_type, response).start
end

Given(/^the SIRI server waits (GetStopMonitoring) request to respond with$/) do |message_type, response|
  SIRIServer.find("default").expect_request(message_type, response)
end

When(/^the SIRI server has received a (GetStopMonitoring) request$/) do |message_type|
  SIRIServer.find("default").wait_request message_type
end

When(/^the SIRI server has received (\d+) (GetStopMonitoring) requests$/) do |count, message_type|
  SIRIServer.find("default").wait_request message_type, count.to_i
end

When(/^I send this SIRI request(?: to the Referential "([^"]*)")?$/) do |referential, request|
  send_siri_request request, referential: referential
end

Then(/^I should receive this SIRI response$/) do |expected_xml|
  save_siri_messages expected: normalized_xml(expected_xml), received: normalized_xml(@last_siri_response), received_raw: @last_siri_response
  expect(normalized_xml(@last_siri_response)).to eq(normalized_xml(expected_xml))
end

Then(/^I should receive a SIRI GetStopMonitoringResponse with$/) do |expected|
  document = REXML::Document.new(@last_siri_response)

  expected_values = {}
  expected.raw.each do |row|
    expected_values[row[0]] = row[1] unless row[2] && row[2] =~ /^TODO/
  end

  actual_values = {}
  expected_values.keys.each do |xpath|
    node = REXML::XPath.first(document, xpath, { "siri" => "http://www.siri.org.uk/siri" })
    xml_value = node.text if node
    actual_values[xpath] = xml_value
  end

  expect(actual_values).to eq(expected_values)
end

When(/^I send a SIRI GetStopMonitoring request with$/) do |attributes|
  default_attributes = {
    "RequestorRef" => "test",
    "MonitoringRef" => "NINOXE:StopPoint:SP:24:LOC"
  }
  attributes = default_attributes.merge(attributes.rows_hash)

  request = %Q{
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>#{attributes['RequestorRef']}</ns2:RequestorRef>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
      </ServiceRequestInfo>

      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
        <ns2:StartTime>2017-01-01T12:00:00.000Z</ns2:StartTime>
        <ns2:MonitoringRef>#{attributes['MonitoringRef']}</ns2:MonitoringRef>
        <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
      </Request>
      <RequestExtension />
    </ns7:GetStopMonitoring>
  </S:Body>
</S:Envelope>
}

  send_siri_request request
end

Then(/^the SIRI server should not have received a (GetStopMonitoring) request$/) do |request|
  received_request = SIRIServer.find("default").received_request?
  expect(received_request).to be_truthy
end

Then(/^the (?:"([^"]*)" )?SIRI server should have received a GetStopMonitoring request with:$/) do |name, attributes|
  name ||= "default"
  last_siri_request = SIRIServer.find(name).requests.last.body

  document = XML::Document.new(last_siri_request)

  expected_values = attributes.rows_hash
  actual_values = document.values(expected_values.keys)

  expect(actual_values).to eq(expected_values)
end
