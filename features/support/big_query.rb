class BigQuery
  include Singleton

  def events
    @events ||= []
  end

  def initialize
    http_server.mount_proc '/' do |request, response|
      serve(request, response)
    end
  end

  PORT = 8349

  def self.url
    "http://localhost:#{PORT}"
  end

  def http_server
    @http_server ||= WEBrick::HTTPServer.new(Port: PORT, Logger: WEBrick::Log.new(File::NULL), AccessLog: [])
  end

  def serve(request, _)
    events << JSON.parse(request.body)
  end

  def started
    @started
  end

  def start
    return if started

	  Thread.start do
      @started = true
		  http_server.start
      @started = false
	  end

    wait { started }
  end

  def stop
    return unless started
    http_server.shutdown
    events.clear

    wait { ! started }
  end

  def clear
    events.clear
  end

  def received_events
    wait { ! events.empty? }
    events
  end

  def wait(&block)
    try_count = 0
    while !block.call
      try_count += 1
      return if try_count > 10

      sleep 0.5
    end
  end

  def self.received_events
    instance.received_events
  end

end

Before do
  BigQuery.instance.start
end

After do
  BigQuery.instance.clear
end
