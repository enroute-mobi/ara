Given(/^the table "([^"]*)" has the following data:$/) do |table_name, datas|
  data_array = datas.raw

  request_string = "INSERT INTO #{table_name} (#{data_array.shift.join(',')}) VALUES"

  data_array.each do |data|
    request_string += "(#{data.join(',')}),"
  end
  request_string.gsub!(/,$/, ';')

  conn = PG.connect dbname: $database
  conn.exec(request_string)
end

When(/^I start Edwig$/) do
  start_edwig()
end