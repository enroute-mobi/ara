require 'gtfs/realtime'

def gtfs_url(attributes = {})
  attributes = {
    referential: 'test',
    path: 'gtfs'
  }.merge(attributes.delete_if { |_,v| v.nil? })

  url_for(attributes)
end
