require 'fileutils'
require 'tmpdir'

$server = 'http://localhost:8081'
$adminToken = "6ceab96a-8d97-4f2a-8d69-32569a38fc64"
$token = "testtoken"

class Ara
  def self.instance
    @ara ||= new
  end

  def root
    @root ||= Pathname.new(File.expand_path(ENV.fetch('ARA_ROOT', Dir.getwd)))
  end

  def tmp_dir
    @tmp_dir ||= Pathname.new(Dir.tmpdir)
  end

  def log_dir
    root_log = root.join('log')
    if root_log.exist?
      root_log
    else
      tmp_dir
    end
  end

  def log_file
    log_dir.join('ara.log')
  end

  def pid_file
    @pid_file ||= tmp_dir.join('ara.pid')
  end

  def config_dir
    root.join('config')
  end

  def initialize
    unless File.directory?("tmp")
      FileUtils.mkdir_p("tmp")
    end
    unless File.directory?("log")
      FileUtils.mkdir_p("log")
    end
  end

  attr_writer :fakeuuid_legacy
  def fakeuuid_legacy?
    @fakeuuid_legacy.nil? ? true : @fakeuuid_legacy
  end

  def environment
    {
      ARA_ROOT: root,
      ARA_CONFIG: config_dir,
      ARA_ENV: 'test',
      ARA_BIGQUERY_PREFIX: 'cucumber',
      ARA_BIGQUERY_TEST: BigQuery.url,
      ARA_FAKEUUID_REAL: !fakeuuid_legacy?
    }
  end

  def command_environment
    environment.map { |k,v| "#{k}=#{v}" }.join(' ')
  end

  def command_executable
    binary_path = root.join('ara')
    if File.exist?(binary_path)
      binary_path
    else
      'go run ara.go'
    end
  end

  def command(arguments)
    "#{command_environment} #{command_executable} #{arguments} >> #{log_file} 2>&1"
  end

  def run(arguments, background: false)
    run_command = command(arguments)
    run_command = "#{run_command} &" if background

    Dir.chdir root do
      system run_command
    end
  end

  def self.load(referential_slug, file)
    Ara.instance.run("load #{file} #{referential_slug}")
  end

  def self.load_content(referential_slug, content)
    Tempfile.open(["ara-import",".csv"]) do |file|
      file.write content
      file.close

      load referential_slug, file.path
    end
  end

  def start
    run "-debug -pidfile=#{pid_file} -testuuid -testclock=20170101-1200 api -listen=localhost:8081", background: true

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

    @started = true
  end

  def started?
    @started
  end

  def self.stop
    return unless @ara
    @ara.stop
    @ara = nil
  end

  def stop
    return unless started?

    pid = IO.read(pid_file)
    Process.kill('KILL',pid.to_i)
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
