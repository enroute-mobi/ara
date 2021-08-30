require 'webrick'
require 'uri'

class GtfsServer

  @@servers = {}
  def self.each(&block)
    @@servers.values.each(&block)
  end

  def self.create(name, url)
    @@servers[name] ||= GtfsServer.new(url)
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
      if ENV["GTFS_DEBUG"]
        puts "Receive GTFS request:"
        puts req
      end
			self.requests << req

		  res.body = self.responses.shift.to_s
		  res.content_type = "application/x-protobuf"
	  end
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

  def expect_request(response)
    if response.is_a? String
      response = GTFS::Realtime::FeedMessage.parse_from_text response
    end

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

end

After do
  GtfsServer.stop
end
