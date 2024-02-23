# module GraphQlHelper

#   def a_graphql_resource_having_attributes(attributes)
#     attributes = attributes.rows_hash if attributes.respond_to?(:rows_hash)

#     attributes_with_human_names = attributes
#     attributes = {}

#     attributes_with_human_names.each do |human_name, value|
#       graphql_name = human_name.parameterize(separator: '_').camelize(:lower)
#       attributes[graphql_name] = value
#     end

#     attribute_matchers = attributes.map do |attribute, value|
#       matcher =
#         case value
#         when %r{^"(.*)"$}
#           $1
#         when %r{^/(.*)/$}
#           match(Regexp.new($1))
#         when %r{^-?\d+$}
#           value.to_i
#         when %r{^-?\d+\.\d+$}
#           a_value_within(0.001).of(value.to_f)
#         when "true"
#           be_truthy
#         when "false"
#           be_falsy
#         when %r{^{.*}$}
#           JSON.parse(value)
#         else
#           eq(value)
#         end
#       [ attribute, matcher ]
#     end.to_h

#     a_hash_including(attribute_matchers)
#   end
# end

# World(GraphQlHelper)

When('I send this GraphQL query to the Referential {string} with token {string}') do |referential, partner_token, query|
  url = URI(url_for referential: referential, path: "graphql")
  query = "{ \"query\": \"#{query.gsub('"','\"').gsub("\n"," ").gsub(/ +/," ")}\" }"

  puts query if ENV['GRAPHQL_DEBUG']
  response = Net::HTTP.post(url, query, "Content-Type" => "application/json", "Authorization" => "Token token=#{partner_token}")
  puts response.body if ENV['GRAPHQL_DEBUG']

  raise "Fail to perform GraphQL query: #{response.value} #{response.body}" unless response.is_a?(Net::HTTPSuccess)

  @graphql_response = JSON.parse(response.body)
end


Then('the GraphQL response should contain one Vehicle with these attributes:') do |attributes|
  expect(@graphql_response["data"]["vehicle"]).to eq(model_attributes attributes)
end

Then('the GraphQL response should contain an updated Vehicle with these attributes:') do |attributes|
  expect(@graphql_response["data"]["updateVehicle"]).to eq(model_attributes attributes)
end

Then('the GraphQL response should contain a Vehicle with these attributes:') do |attributes|
  collection = @graphql_response["data"]["vehicles"]

  expect(collection).to include(model_attributes attributes)
end
