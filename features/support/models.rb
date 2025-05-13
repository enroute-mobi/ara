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
