require 'webrick'
require 'uri'
require 'base64'

class OAUTHServer
  @@servers = {}
  def self.each(&block)
    @@servers.values.each(&block)
  end

  def self.create(name, url, options)
    @@servers[name] ||= OAUTHServer.new(url, options)
  end

  def self.find(name)
    @@servers[name]
  end

  def self.stop
    each(&:stop)
    @@servers.clear
  end

  attr_accessor :url, :started

  def initialize(url, options)
    @url = url
    @uri = URI.parse(@url)
    @options = options
    @http_server = WEBrick::HTTPServer.new(
      Port: @uri.port,
      Logger: WEBrick::Log.new(File::NULL),
      AccessLog: []
    )

    @http_server.mount_proc @uri.path do |req, res|
      authorization_header = Base64::decode64(req.header['authorization'][0].split[1])

      if authorization_header == expected_bearer && req.body == expected_body
        payload = {'access_token': @options['access_token'],
                'token_type': 'Bearer',
                'expires_in': 3600,
                'id_token': 'eyJhSUzI1Ni---Y0ZDQEifQ'}

        res.content_type = 'application/json'
        res.body = JSON.dump(payload)

        # we have to warn SIRI server that a token is needed
        SIRIServer.authorized_tokens << @options['access_token']
      else
        res.status = 401
        res.body = "Wrong Bearer for authotization, expected: #{expect_bearer}, got: #{authorizatio_header}"
      end
    end
  end

  def expected_bearer
    "#{@options['client_id']}:#{@options['client_secret']}"
  end

  def expected_scopes
    @options['scopes']&.split(',')&.join('+')&.to_s
  end

  def expected_body
    body = 'grant_type=client_credentials'
    if !expected_scopes.nil?
      body = body + '&' + expected_scopes
    end

    body
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
  OAUTHServer.stop
end
