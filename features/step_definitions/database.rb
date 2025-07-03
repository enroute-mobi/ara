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

When(/^I start Ara$/) do
  TestAra.instance.start
end

Then('the table {string} has rows with the following values:') do |table_name, values|
  row_values = values.hashes
  result = @connection.exec("select * from #{table_name};")
  expect(result.to_a).to include(*row_values.map { |r| a_hash_including(r) })
end

Then('the table {string} has a row with the following values:') do |table_name, values|
  row_values = values.rows_hash
  result = @connection.exec("select * from #{table_name};")
  expect(result.to_a).to include(a_hash_including(row_values))
end

Then('the table {string} has no row with the following values:') do |table_name, values|
  row_values = values.rows_hash
  result = @connection.exec("select * from #{table_name};")
  expect(result.to_a).to_not include(a_hash_including(row_values))
end
