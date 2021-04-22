require 'fileutils'

$server = 'http://localhost:8081'
$adminToken = "6ceab96a-8d97-4f2a-8d69-32569a38fc64"
$token = "testtoken"

class Ara
  def self.instance
    @ara ||= new
  end

  attr_writer :fakeuuid_legacy
  def fakeuuid_legacy?
    @fakeuuid_legacy.nil? ? true : @fakeuuid_legacy
  end

  def start
    unless File.directory?("tmp")
      FileUtils.mkdir_p("tmp")
    end
    unless File.directory?("log")
      FileUtils.mkdir_p("log")
    end

    ara_command = "ARA_ROOT=#{Dir.getwd} ARA_CONFIG=#{Dir.getwd}/config ARA_ENV=test ARA_BIGQUERY_PREFIX=cucumber ARA_BIGQUERY_TEST=#{BigQuery.url} ARA_FAKEUUID_REAL=#{!fakeuuid_legacy?} go run ara.go -debug -pidfile=tmp/pid -testuuid -testclock=20170101-1200 api -listen=localhost:8081 >> log/ara.log 2>&1 &"
    system ara_command

    time_limit = Time.now + 30
    while
      sleep 0.5

      begin
        response = RestClient::Request.execute(method: :get, url: "#{$server}/_status", timeout: 1)
        break if response.code == 200 && JSON.parse(response.body)["status"] == "ok"
      rescue StandardError

      end

      raise "Timeout" if Time.now > time_limit
    end
  end

  def self.stop
    pid = IO.read("tmp/pid")
    Process.kill('KILL',pid.to_i)

    @ara = nil
  end

end

Before('@database') do
  Ara.instance.fakeuuid_legacy = false
end

Before('not @nostart') do
  Ara.instance.start
end

After do
  Ara.stop
end
