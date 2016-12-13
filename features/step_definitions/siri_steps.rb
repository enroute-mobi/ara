# Given(/^we send a checkstatus request with body "([^"]+)"$/) do |filename|
#   system "rm features/testdata/response.xml"
#   system "curl -s -H 'Content-Type: text/xml' #{$server}/default/siri --data @features/testdata/#{filename} > features/testdata/response.xml"
# end

# Then(/^we should recieve a checkstatus response with body "([^"]+)"$/) do |filename|
#   expect(IO.read("features/testdata/response.xml")).to eq(IO.read("features/testdata/#{filename}"))
# end
