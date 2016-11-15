require 'fileutils'

$server = 'http://localhost:8081'

Before('@server') do
  unless File.directory?("tmp")
    FileUtils.mkdir_p("tmp")
  end
  system "go run edwig.go -pidfile=tmp/pid -testuuid -testclock=20170101-1200 api -listen=localhost:8081 &"

  time_limit = Time.now + 10
  begin
    sleep 2
    system "go run edwig.go check #{$server}/default/siri"
    raise "Timeout" if Time.now > time_limit
  end until $?.exitstatus == 0
end

After('@server') do
  pid = IO.read("tmp/pid")
  Process.kill('KILL',pid.to_i)
end
