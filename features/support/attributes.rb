def model_attributes(table)
  attributes = table.rows_hash.dup

  attributes.dup.each do |key, value|
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

    case value
    when /\A\d+\Z/
      # Convert integer
      attributes[key] = value.to_i
    end
  end

  if objectids = (attributes["ObjectIds"] || attributes["ObjectIDs"])
    attributes["ObjectIDs"] = JSON.parse("{ #{objectids} }")
  end

  attributes
end

def api_attributes(json)
  attributes = JSON.parse(json)

  objectids = attributes["ObjectIDs"]
  if Array === objectids
    attributes["ObjectIDs"] = Hash[objectids.map { |objectid| [objectid["Kind"], objectid["Value"]] }]
  end

  attributes
end
