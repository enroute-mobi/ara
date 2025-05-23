def find_referential(slug)
  slug ||= 'test'
  referential = TestAra.instance.server.referentials.find(slug)
  expect(referential).not_to be_nil

  referential
end

def find_model(slug, model, value)
  referential = find_referential(slug)
  referential.send(model).find(value)
end

def referential_models(slug, model)
  referential = find_referential(slug)
  referential.send(model).all
end

ParameterType(
  name: 'ara_resource',
  regexp: Regexp.new(
    %w[StopArea Line Vehicle VehicleJourney Operator StopVisit ScheduledStopVisit].join('|')
  ),
  transformer: ->(s) { s }
)

Given("a {ara_resource} is created with the following attributes:") do |collection, attributes|
  step "a #{collection} is created in Referential \"test\" with the following attributes:", attributes
end

Given("a {ara_resource} is created in Referential {string} with the following attributes:") do |collection, slug, attributes|
  referential = find_referential(slug)
  model = referential.send("#{collection.underscore}s").create(model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym })

  raise "Cannot create #{collection}: #{model.errors}" unless model.save
end

When("an {ara_resource} exists with the following attributes:") do |collection, attributes|
  step "a #{collection} is created with the following attributes:", attributes
end

When("a {ara_resource} exists with the following attributes:") do |collection, attributes|
  step "a #{collection} is created with the following attributes:", attributes
end

Then("one {ara_resource} has the following attributes:") do |collection, attributes|
  step "one #{collection} in Referential \"test\" has the following attributes:", attributes
end

Then("one {ara_resource} in Referential {string} has the following attributes:") do |collection, slug, attributes|
  referential = find_referential(slug)

  check_attributes(referential, "#{collection.underscore}s", attributes)
end

Then("the {ara_resource} {string} has the following attributes:") do |collection, code_or_id, attributes|
  step "the #{collection} \"#{code_or_id}\" in Referential \"test\" has the following attributes:", attributes
end

Then("the {ara_resource} {string} in Referential {string} has the following attributes:") do |collection, code_or_id, slug, attributes|
  model = find_model(slug, "#{collection.underscore}s", code_or_id)

  matcher_attributes(attributes, model)
end

When('the {ara_resource} {string}:{string} is destroyed') do |ara_resource, code_space, value|
  model = find_model('test', "#{ara_resource.underscore}s", "#{code_space}:#{value}")

  raise "Cannot destroy #{ara_model}: #{model.errors}" unless model.destroy
end

Then("a {ara_resource} {string}:{string} should exist") do |collection, code_space, value|
  step "a #{collection} \"#{code_space}\":\"#{value}\" should exist in Referential \"test\""
end

Then("a {ara_resource} {string}:{string} should exist in Referential {string}") do |collection, code_space, value, slug|
  model = find_model(slug, "#{collection.underscore}s", "#{code_space}:#{value}")
  expect(model).not_to be_nil
end

Then("a {ara_resource} {string}:{string} should not exist") do |collection, code_space, value|
  step "a #{collection} \"#{code_space}\":\"#{value}\" should not exist in Referential \"test\""
end

Then("a {ara_resource} {string}:{string} should not exist in Referential {string}") do |collection, code_space, value, slug|
  model = find_model(slug, "#{collection.underscore}s", "#{code_space}:#{value}")
  expect(model).to be_nil
end

Then("the {ara_resource} {string} is edited with the following attributes:") do |collection, code_or_id, attributes|
  step "the #{collection} \"#{code_or_id}\" is edited in Referential \"test\" with the following attributes:", attributes
end

Then("the {ara_resource} {string} is edited in Referential {string} with the following attributes:") do |collection, code_or_id, slug, attributes|
  model = find_model(slug, "#{collection.underscore}s", code_or_id)
  model.model_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  raise "Cannot update #{collection}" unless model.save
end
