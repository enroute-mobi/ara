Before('@server') do
  system "go run edwig.go -testuuid -testclock=20170101-1200 api &"
  # $server = $?.pid

  time_limit = Time.now + 10
  begin
    sleep 2
    system "go run edwig.go check http://localhost:8080/siri"
    raise "Timeout" if Time.now > time_limit
  end until $?.exitstatus == 0
end

After('@server') do
  system "killall edwig"
  # ps = `ps`
  # puts "$?.pid = #{$server}"
  # puts "ps output:"
  # puts "#{ps}"
  # Process.kill('KILL',$server)
end