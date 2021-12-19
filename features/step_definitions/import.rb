Given('I import these models in the referential {string}:') do |referential_slug, content|
  Ara.load_content referential_slug, content
end

Then('I can import these models in the referential {string}:') do |referential_slug, content|
  expect(Ara.load_content(referential_slug, content)).to be_truthy
end
