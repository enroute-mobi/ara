Given(/^the table "([^"]*)" has the following data:$/) do |table_name, datas|
  data_array = datas.raw

  request_string = "INSERT INTO #{table_name} (#{data_array.shift.join(',')}) VALUES"

  data_array.each do |data|
    request_string += "(#{data.join(',')}),"
  end
  request_string.gsub!(/,$/, ';')

  # conn = PG.connect dbname: $database, user: ENV["POSTGRESQL_ENV_POSTGRES_USER"], password: ENV["POSTGRESQL_ENV_POSTGRES_PASSWORD"]
  @connection.exec(request_string)
end

When(/^I start Edwig$/) do
  start_edwig()
end
