Given(/^an Operator exists (?:in Referential "([^"]+)" )?with the following attributes:$/) do |slug, operator|
  referential = find_referential(slug)
  operator = referential.operators.create(model_attributes(operator).transform_keys { |key| key.to_s.underscore })

  raise 'Cannot create operator' unless operator.save
end
