def audit_attributes(attributes)
  attributes = attributes.rows_hash

  # Rewrite values like numbers, nil, etc
  attributes.map do |key, value|
    value =
      case value
      when /^\d+$/
        value.to_i
      when "nil"
        nil
      when "<empty>"
        ""
      else
        value
      end
    [ key, value ]
  end.to_h
end

Then('an audit event should exist with these attributes:') do |attributes|
  expect(BigQuery.received_events).to include(a_hash_including(audit_attributes(attributes)))
end

Then('an audit event should not exist with these attributes:') do |attributes|
  expect(BigQuery.received_events).to_not include(a_hash_including(audit_attributes(attributes)))
end
