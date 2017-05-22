require 'webrick'
require 'uri'
require 'optparse'
require 'securerandom'
require 'time'
require 'json'

options = {}
OptionParser.new do |parser|
  parser.banner = "Usage: mock-consumer.rb [options]"

  parser.on("--listen=URL", "Mock server address") do |url|
    unless url[/\Ahttps?:\/\//]
      url = "http://#{url}"
    end
    options[:listen] = URI.parse(url)
  end

  parser.on("--logstash=URL", "Logstash address") do |url|
    unless url[/\Ahttps?:\/\//]
      url = "http://#{url}"
    end
    options[:logstash] = URI.parse(url)
  end
end.parse!

if options.length < 2
  puts "Wrong number of arguments"
  exit
end

startedTime = Time.now
http_server = WEBrick::HTTPServer.new(Port: options[:listen].port, Logger: WEBrick::Log.new(File::NULL), AccessLog: [])

http_server.mount_proc options[:listen].path do |req, res|
  log = {request: req}
  if req.body =~ /CheckStatus/
    res.content_type = "text/xml"
    res.body = %Q{<?xml version='1.0' encoding='utf-8'?>
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
        <ns3:ResponseTimestamp>#{Time.now.strftime("%FT%T%:z")}</ns3:ResponseTimestamp>
        <ns3:ProducerRef>RatpDev</ns3:ProducerRef>
        <ns3:ResponseMessageIdentifier>#{SecureRandom.uuid}</ns3:ResponseMessageIdentifier>
        <ns3:RequestMessageRef>#{req.body.match(/MessageIdentifier>([^<\s]+)<\//)[1]}</ns3:RequestMessageRef>
      </CheckStatusAnswerInfo>
      <Answer>
        <ns3:Status>true</ns3:Status>
        <ns3:ServiceStartedTime>#{startedTime.strftime("%FT%T%:z")}</ns3:ServiceStartedTime>
      </Answer>
      <AnswerExtension />
    </ns8:CheckStatusResponse>
  </S:Body>
</S:Envelope>}
  log[:response] = res.body
  else
    res.status = 200
  end
  # Log to Logstash
  logstash = TCPSocket.new options[:logstash].host, options[:logstash].port
  logstash.write(log)
  logstash.close
end

http_server.start