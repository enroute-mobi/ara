def model_attributes(table)
  attributes = table.rows_hash.dup

  attributes.dup.each do |key, value|
    case value
    when /\A\d+\Z/
      # Convert integer
      attributes[key] = value.to_i
    when /\A(true|false)\Z/
      # Convert boolean
      attributes[key] = (value == "true")
    when /\A\[.+\]\Z/
      # Convert Array
      attributes[key] = JSON.parse(value)
    end

    # Transform
    #  | Schedule[aimed]#Arrival   | 2017-01-01T13:00:00.000Z          |
    #  | Schedule[aimed]#Departure | 2017-01-01T13:02:00.000Z          |
    # into
    # "Schedules" => [
    #   {"Kind"=>"aimed", "ArrivalTime"=>"2017-01-01T13:00:00.000Z", "DepartureTime"=>"2017-01-01T13:02:00.000Z"}
    # ]
    if key =~ /Schedule\[(aimed|expected|actual)\]#(Arrival|Departure)/
      schedule_type = $1
      attribute = $2

      attributes["Schedules"] ||= []

      schedule = attributes["Schedules"].find { |s| s["Kind"] == schedule_type }
      unless schedule
        schedule = { "Kind" => schedule_type }
        attributes["Schedules"] << schedule
      end

      schedule["#{attribute}Time"] = value

      attributes.delete key
    end

    if key =~ /Messages\[(\d+)\]#(\S+)/
      message_number = $1.to_i
      attribute = $2

      attributes["Messages"] ||= []

      until attributes["Messages"].length >= message_number+1
        attributes["Messages"] << {}
      end
      message = attributes["Messages"][message_number]

      message[attribute] = value
      attributes.delete key
    end

    if key =~ /Attribute\[([^\]]+)\]/
      name = $1
      attributes["Attributes"] ||= {}
      attributes["Attributes"][name] = value
      attributes.delete key
    end

    # Situation References are an array of Reference
    # Format: | Reference[0] | Kind:ObjectId |
    if key =~ /References\[(\d+)\]/
      reference_number = $1.to_i

      attributes["References"] ||= []

      until attributes["References"].length >= reference_number+1
        attributes["References"] << {}
      end

      kind, objectid = value.split(":",2)
      attributes["References"][reference_number] = { "Type" => kind, "ObjectId" => JSON.parse(objectid) }
      attributes.delete key
    end

    if key =~ /Reference\[([^\]]+)\]#(ObjectID|Id)/
      name = $1
      attribute = $2
      attributes["References"] ||= {}

      if attribute == "ObjectID"
        value = JSON.parse("{ #{value} }")
      end

      attributes["References"][name] = { attribute => value }
      attributes.delete key
    end

    if key =~ /ReferenceArray(\d+)/
      name = $1
      attribute = $2

      attributes["References"] ||= []

      values = value.split(',')
      attributes["References"][$1.to_i] = {
        "Type" => values[0],
        "ObjectID" => JSON.parse("{ #{values[1]} }")
      }

      attributes.delete key
    end
  end

  if objectids = (attributes["ObjectIDs"] || attributes["ObjectIDs"])
    attributes["ObjectIDs"] = JSON.parse("{ #{objectids} }")
  end

  if settings = attributes["Settings"]
    attributes["Settings"] = JSON.parse(settings)
  end

  attributes
end

def api_attributes(json)
  JSON.parse(json)
  # puts json.inspect
  # attributes = (String === json ? JSON.parse(json) : json)
  # puts attributes.inspect

  # if Array === attributes
  #   return attributes.map { |item_attributes| api_attributes(item_attributes) }
  # end

  # objectids = attributes["ObjectIDs"]
  # if Array === objectids
  #   attributes["ObjectIDs"] = Hash[objectids.map { |objectid| [objectid["Kind"], objectid["Value"]] }]
  # end

  # attributes
end

def has_attributes(response_array, attributes)
  parsed_attributes = model_attributes(attributes)

  objectid_kind = parsed_attributes["ObjectIDs"].keys.first
  objectid_value = parsed_attributes["ObjectIDs"][objectid_kind]

  found_value = response_array.find{|a| a["ObjectIDs"][objectid_kind] == objectid_value}

  expect(found_value).not_to be_nil

  parsed_attributes.delete("ObjectIDs")
  expect(found_value).to include(parsed_attributes)
end
