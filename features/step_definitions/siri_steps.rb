# Given(/^a SIRI HTTP server$/) do
#   system "go run edwig.go -testuuid -testclock=20170101-1200 api &"
# end

Given(/^we send a checkstatus request with body "([^"]+)"$/) do |filename|
  system "rm features/testdata/response.xml"
  system "curl -s -H 'Content-Type: text/xml' http://localhost:8080/default/siri --data @features/testdata/#{filename} > features/testdata/response.xml"
end

Then(/^we should recieve a checkstatus response with body "([^"]+)"$/) do |filename|
  expect(IO.read("features/testdata/response.xml")).to eq(IO.read("features/testdata/#{filename}"))
end
