require 'fileutils'

$server = 'http://localhost:8081'

Before('@server') do
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

After('@server') do
  pid = IO.read("tmp/pid")
  Process.kill('KILL',pid.to_i)
end
