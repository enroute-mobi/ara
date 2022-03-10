require 'csv'
require 'net/http'

def import_path(referential_slug)
  url_for(path: "#{referential_slug}/import")
end

def prepare_csv(doc_string)
  tmp = Tempfile.new(['import_test', '.csv'])
  tmp.write(doc_string)
  tmp.close
  tmp
end

def send_to_ara(file, referential, token)
  payload = { data: file.open, request: { force: true }.to_json }

  @response = begin
                RestClient.post(
                  import_path(referential),
                  payload,
                  { Authorization: "Token token=#{token}" }
                )
              rescue => e
                e
              end
end

Given('I import these models in the referential {string}:') do |referential_slug, content|
  Ara.load_content referential_slug, content
end

Then('I can import these models in the referential {string}:') do |referential_slug, content|
  expect(Ara.load_content(referential_slug, content)).to be_truthy
end

When('I import in the referential {string} with the token {string} these models:') do |referential_slug, token, doc_string|
  send_to_ara(prepare_csv(doc_string), referential_slug, token)
end

Then('the import should be successful') do
  expect(@response.code).to eq(200)
  expect(JSON.parse(@response.body)).to include('Errors' => {})
end

Then('the import must fail with an unauthorized status') do
  expect(@response.http_code).to eq(401)
end
