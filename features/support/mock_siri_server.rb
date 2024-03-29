require 'webrick'
require 'uri'

class SIRIServer

  @@servers = {}
  @@authorized_tokens = []

  def self.each(&block)
    @@servers.values.each(&block)
  end

  def self.create(name, url, envelope = 'SOAP')
    @@servers[name] ||= SIRIServer.new(url, envelope)
  end

  def self.find(name)
    @@servers[name]
  end

  def self.stop
    each(&:stop)
    @@servers.clear
  end

  def self.authorized_tokens
    @@authorized_tokens
  end

  attr_accessor :url, :port, :path, :requests, :responses, :started, :service_started_at

  SERVICE_STARTED_AT = '2016-09-22T03:30:32.000+02:00'.freeze

  def initialize(url, envelope)
    @url = url
    @uri = URI.parse(@url)
    @requests = []
    @responses = []
    @envelope = envelope
    @http_server = WEBrick::HTTPServer.new(Port: @uri.port, Logger: WEBrick::Log.new(File::NULL), AccessLog: [])
    @service_started_at = SERVICE_STARTED_AT

    proceed_request
  end

  def proceed_request
    @http_server.mount_proc @uri.path do |req, res|
      if need_authorization?(req) && !authorize?(req)
        res.status = 401
        res.body = 'Unauthorized request'
      else
        Siri::Validator.create(req.body,"server request").log

        request_message_identifiers = req.body.scan(/MessageIdentifier>(.*)</).flatten
        if checkstatus_response_standard?(req)
          puts "[#{url}] Response to CheckStatus request" if ENV["SIRI_DEBUG"]
          case @envelope
          when 'SOAP'
            res.body = soap_checkstatus_response(request_message_identifiers)
          when 'raw'
            res.body = raw_checkstatus_response(request_message_identifiers)
          else
            raise "Unknown envelope #{@envelope}"
          end
        else
          puts "[#{url}] Receive SIRI request:\n#{req}" if ENV["SIRI_DEBUG"]
          requests << req

          request_body = @responses.shift
          request_body.gsub!("{RequestMessageRef}", request_message_identifiers.first)
          request_body.gsub!("{LastRequestMessageRef}", request_message_identifiers.last)

          res.body = request_body
          puts "[#{url}] Send SIRI response:\n#{res.body}" if ENV["SIRI_DEBUG"]
        end

        Siri::Validator.create(res.body,"server response").log
      end
      res.content_type = 'text/xml'
    end
  end

  def need_authorization?(req)
    !req.header['authorization'].empty?
  end

  def authorize?(req)
    SIRIServer.authorized_tokens.include?(req.header['authorization'].first.split[1])
  end

  def expect_request(type, response)
    @responses << response
    self
  end

  def wait_request(type, count = 1)
    try_count = 0
    while requests.count < count
      try_count += 1
      raise "Received #{requests.count} request" if try_count > 10

      sleep 0.5
    end
  end

  def checkstatus_response_standard?(req)
    req.body =~ /CheckStatus/ && @responses.first !~ /CheckStatus/
  end

  def raw_checkstatus_response(request_message_identifiers)
    %(<?xml version="1.0" encoding="UTF-8"?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<CheckStatusResponse>
    <ResponseTimestamp>2016-09-22T07:58:34.000+02:00</ResponseTimestamp>
    <ProducerRef>NINOXE:default</ProducerRef>
    <ResponseMessageIdentifier>c464f588-5128-46c8-ac3f-8b8a465692ab</ResponseMessageIdentifier>
    <RequestMessageRef>#{request_message_identifiers.first}</RequestMessageRef>
    <Status>true</Status>
    <ServiceStartedTime>#{@service_started_at}</ServiceStartedTime>
</CheckStatusResponse>
</Siri>)
  end

  def soap_checkstatus_response(request_message_identifiers)
    %Q(<?xml version='1.0' encoding='utf-8'?>
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
        <ns3:Address>#{@url}</ns3:Address>
        <ns3:ResponseMessageIdentifier>c464f588-5128-46c8-ac3f-8b8a465692ab</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>#{request_message_identifiers.first}</ns3:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <ns3:Status>true</ns3:Status>
        <ns3:ServiceStartedTime>#{@service_started_at}</ns3:ServiceStartedTime>
      </Answer>
      <AnswerExtension/>
    </ns8:CheckStatusResponse>
  </S:Body>
</S:Envelope>)
  end

  def received_request?
    !requests.empty?
  end

  def received_requests?(count = 1)
    requests.length == count
  end

  def received_specific_requests?(name, count = 1)
    requests.count do |r|
      document = REXML::Document.new(r.body)
      document.root.find_first_recursive {|node| node.name == name }
    end == count
  end

  def start
    return if started

    self.started = true
    Thread.start do
      @http_server.start
    end

    self
  end

  def stop
    @http_server.shutdown
    self.started = false
  end
end

After do
  SIRIServer.stop
end
