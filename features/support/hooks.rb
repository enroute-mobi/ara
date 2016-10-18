require 'fileutils'

Before('@server') do
  unless File.directory?("tmp")
    FileUtils.mkdir_p("tmp")
  end
  system "go run edwig.go -pidfile=tmp/pid -testuuid -testclock=20170101-1200 api &"

  time_limit = Time.now + 10
  begin
    sleep 2
    system "go run edwig.go check http://localhost:8080/siri"
    raise "Timeout" if Time.now > time_limit
  end until $?.exitstatus == 0
end

After('@server') do
  pid = IO.read("tmp/pid")
  Process.kill('KILL',pid.to_i)
end