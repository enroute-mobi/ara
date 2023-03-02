def url_for(attributes = {})
  a = {
    server: $server
  }.merge(attributes.delete_if { |k,v| v.nil? })

  url_parts = [ a[:server], a[:referential], a[:path] ]
  url_parts.compact.join('/').tap do |url|
    # puts a.inspect
    # puts url
  end
end

def url_for_model(attributes = {})
  raise "No specified resource" unless attributes.has_key? :resource

  attributes = {
    referential: 'test'
  }.merge(attributes.delete_if { |k,v| v.nil? })

  if attributes[:model] == 'subscriptions'
    path = [ "#{attributes[:resource]}s", attributes[:partner_name], attributes[:model]].compact.join('/')
  else
    path = [ "#{attributes[:resource]}s", attributes[:id] ].compact.join('/')
  end
  url_for(attributes.merge(path: path))
end
