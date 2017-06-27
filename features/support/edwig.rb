require 'fileutils'

$server = 'http://localhost:8081'
$adminToken = "6ceab96a-8d97-4f2a-8d69-32569a38fc64"
$token = "testtoken"

def start_edwig
  unless File.directory?("tmp")
    FileUtils.mkdir_p("tmp")
  end
  unless File.directory?("log")
    FileUtils.mkdir_p("log")
  end
  system "EDWIG_ENV=test go run edwig.go -debug -pidfile=tmp/pid -testuuid -testclock=20170101-1200 api -listen=localhost:8081 >> log/edwig.log 2>&1 &"

  time_limit = Time.now + 30
  while
    sleep 0.5

    begin
      response = RestClient::Request.execute(method: :get, url: "#{$server}/_status", timeout: 1, :headers => {:Authorization => 'Token token=6ceab96a-8d97-4f2a-8d69-32569a38fc64'})
      break if response.code == 200 && JSON.parse(response.body)["status"] == "ok"
    rescue Exception # => e
      # puts e.inspect
    end

    raise "Timeout" if Time.now > time_limit
  end
end

Before('~@database') do
  start_edwig()
end

After do
  pid = IO.read("tmp/pid")
  Process.kill('KILL',pid.to_i)
end
