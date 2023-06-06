require 'webrick'
require 'uri'

class SIRILiteServer
  @@servers = {}
  def self.each(&block)
    @@servers.values.each(&block)
  end

  def self.create(name, url)
    @@servers[name] ||= SIRILiteServer.new(url)
  end

  def self.find(name)
    @@servers[name]
  end

  def self.stop
    each(&:stop)
    @@servers.clear
  end

  attr_accessor :url, :port, :path, :requests, :responses, :started

  def initialize(url)
    @url = url
    @requests = []
    @responses = []

    uri = URI.parse(url)
    @http_server = WEBrick::HTTPServer.new(Port: uri.port, Logger: WEBrick::Log.new(File::NULL), AccessLog: [])

    @http_server.mount_proc uri.path do |req, res|
      if ENV['SIRI_DEBUG']
        puts 'Receive SIRI request:'
        puts req
      end
      requests << req

      res.body = responses.shift.to_s
      res.content_type = 'application/json'
      res.status = 200
    end
  end

  def expect_request(_type, response)
    @responses << response
    self
  end

  def wait_request(count = 1)
    try_count = 0
    while requests.count < count
      try_count += 1
      raise "Received #{requests.count} request" if try_count > 10

      sleep 0.5
    end
  end

  def received_request?
    !requests.empty?
  end

  def received_requests?(count = 1)
    requests.length == count
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
  SIRILiteServer.stop
end
