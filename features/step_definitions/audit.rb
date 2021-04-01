def an_audit_event_with_attributes(attributes)
  attributes = attributes.rows_hash if attributes.respond_to?(:rows_hash)

  # Transform specified values into integer, nil, regexp, etc
  attribute_matchers = attributes.map do |attribute, value|
      matcher =
        case value
        when "nil"
          nil
        when "<empty>"
          ""
        when /^\d+$/
          value.to_i
        when %r{^/(.*)/$}
          definition = $1
          definition.gsub!("{uuid}","\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b")
          definition.gsub!("{test-uuid}","\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{1,4}-\\b[0-9a-f]{12}\\b")
          match(Regexp.new(definition))
        when %r{^\[.*\]$}
          eval(value)
        else
          value
        end
      [ attribute, matcher ]
  end.to_h

  time_reference = Time.utc(2017,1,1,12)
  attribute_matchers["Timestamp"] ||= satisfy("near #{time_reference}") { |t| Time.parse(t) - time_reference < 300 }

  a_hash_including(attribute_matchers)
end

Then('an audit event should exist with these attributes:') do |attributes|
  expect(BigQuery.received_events).to include(an_audit_event_with_attributes(attributes))
end

Then('an audit event should not exist with these attributes:') do |attributes|
  expect(BigQuery.received_events).to_not include(an_audit_event_with_attributes(attributes))
end
