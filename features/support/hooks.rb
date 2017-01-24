require 'fileutils'
require 'rest-client'
require 'json'

$server = 'http://localhost:8081'

def url_for(attributes = {})
  a = {
    server: $server
  }.merge(attributes.delete_if { |k,v| v.nil? })

  url_parts = [ a[:server], a[:referential], a[:path] ]
  url_parts.compact.join('/').tap do |url|
    # puts a.inspect
    # puts url
  end
end

def url_for_model(attributes = {})
  raise "No specified resource" unless attributes.has_key? :resource

  attributes = {
    referential: 'test'
  }.merge(attributes.delete_if { |k,v| v.nil? })

  path = [ "#{attributes[:resource]}s", attributes[:id] ].compact.join('/')
  url_for(attributes.merge(path: path))
end

def model_attributes table
  attributes = table.rows_hash
  if attributes["ObjectIds"]
    attributes["ObjectIds"] = JSON.parse("{#{attributes["ObjectIds"]}}")
  end
  attributes
end

Before do
  unless File.directory?("tmp")
    FileUtils.mkdir_p("tmp")
  end
  unless File.directory?("log")
    FileUtils.mkdir_p("log")
  end
  system "EDWIG_ENV=test go run edwig.go -debug -pidfile=tmp/pid -testuuid -testclock=20170101-1200 api -listen=localhost:8081 >> log/edwig.log 2>&1 &"

  time_limit = Time.now + 30
  begin
    sleep 2
    system "go run edwig.go check #{$server}/default/siri > /dev/null 2>&1"
    raise "Timeout" if Time.now > time_limit
  end until $?.exitstatus == 0
end

After do
  pid = IO.read("tmp/pid")
  Process.kill('KILL',pid.to_i)
end
