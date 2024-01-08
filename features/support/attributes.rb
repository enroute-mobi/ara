def model_attributes(table)
  attributes = table.rows_hash.dup

  attributes.dup.each do |key, value|
    case value
    when /\A"\d+.\d+"\Z/
      # Don't convert integer between quotes
      attributes[key] = value[1..-2]
    when /\A\d+\Z/
      # Convert integer
      attributes[key] = value.to_i
    when /\A\d+\.\d+\Z/
      # Convert float
      attributes[key] = value.to_f
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

    # Transform
    #  | Origin[partner]  | true  |
    #  | Origin[partner2] | false |
    # into
    # "Origins" => {"partner"=>true, "partner2"=>false}
    if key =~ /Origin\[([^\]]+)\]/
      partner = $1

      attributes["Origins"] ||= {}

      attributes["Origins"][$1] = value == "true" unless attributes["Origins"].key?($1)

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

    if key =~ /Description\[([^\]]+)\]/
      name = $1
      attributes["Description"] ||= {}
      attributes["Description"][name] = value
      attributes.delete key
    end

    if key =~ /Summary\[([^\]]+)\]/
      name = $1
      attributes["Summary"] ||= {}
      attributes["Summary"][name] = value
      attributes.delete key
    end

    # Situation ValidityPeriods is an array of TimeRange
    # Format: | ValidityPeriods[0]#StartTime | 2017-01-01T13:00:00.000Z |
    #         | ValidityPeriods[0]#EndTime   | 2017-01-02T15:00:00.000Z |
    #
    # Situation PublicationWindows Periods is an array of TimeRange
    # Format: | PublicationWindows[0]#StartTime | 2017-01-01T13:00:00.000Z |
    #         | PublicationWIndows[0]#EndTime   | 2017-01-02T15:00:00.000Z |
    # Situation Periods is an array of TimeRange
    # Format: | Periods[0]#StartTime | 2017-01-01T13:00:00.000Z |
    #         | Periods[0]#EndTime   | 2017-01-02T15:00:00.000Z |
    if key =~ /(ValidityPeriods|PublicationWindows|Periods)\[(\d+)\]#(\S+)/
      raw_attribute = Regexp.last_match(0).to_s
      name = Regexp.last_match(1).to_s
      period_number = Regexp.last_match(2).to_i
      attribute = Regexp.last_match(3).to_s

      attributes[name] ||= []

      until attributes[name].length >= period_number + 1
        attributes[name] << {}
      end

      attributes[name][period_number][attribute.to_s] = value

      attributes.delete raw_attribute
    end

    # Situation References are an array of Reference
    # Format: | Reference[0] | Kind:Code |
    if key =~ /References\[(\d+)\]/
      reference_number = $1.to_i

      attributes["References"] ||= []

      until attributes["References"].length >= reference_number+1
        attributes["References"] << {}
      end

      kind, code = value.split(":",2)
      attributes["References"][reference_number] = { "Type" => kind, "Code" => JSON.parse(code) }
      attributes.delete key
    end

    if key =~ /Reference\[([^\]]+)\]#(Code|Id)/
      name = $1
      attribute = $2
      attributes["References"] ||= {}

      if attribute == "Code"
        value = JSON.parse("{ #{value} }")
      end

      attributes["References"][name] = { attribute => value }
      attributes.delete key
    end

    if key =~ %r{^(Affects\[([^\]]+)\])(/(AffectedDestinations|AffectedSections|AffectedRoutes)\[(\d+)])?(/((FirstStop|LastStop)|StopAreaId|RouteRef))?}
      raw_attribute = Regexp.last_match(0).to_s
      attribute = Regexp.last_match(2).to_s
      subaffect = Regexp.last_match(4).to_s
      index = Regexp.last_match(5).to_i
      stop_type = Regexp.last_match(7).to_s

      attributes['Affects'] ||= []
      case attribute
      when 'StopArea', 'Line'
        attributes['Affects'] << {
          'Type' => attribute,
          "#{attribute}Id" => value
        }
      else
        model, id = attribute.split('=')
        attributes['Affects'].map do |affect|
          next if affect['Type'] != model && affect["#{model}Id"] != id

          affect[subaffect] ||= []
          affect[subaffect][index] ||= {}
          affect[subaffect][index][stop_type] = value
        end
      end
      attributes.delete raw_attribute
    end
    
    if key =~ /ReferenceArray\[(\d+)\]/
      name = $1
      attribute = $2

      attributes["References"] ||= []

      values = value.split(',')
      attributes["References"][$1.to_i] = {
        "Type" => values[0],
        "Code" => JSON.parse("{ #{values[1]} }")
      }

      attributes.delete key
    end
  end

  if codes = (attributes["Codes"] || attributes["Codes"])
    attributes["Codes"] = JSON.parse("{ #{codes} }")
  end

  if settings = attributes["Settings"]
    attributes["Settings"] = JSON.parse(settings)
  end

  attributes["Schedules"].sort_by!{ |s| s["Kind"] } if attributes["Schedules"]

  attributes
end

def api_attributes(json)
  # puts json.inspect
  attributes = (String === json ? JSON.parse(json) : json)
  # puts attributes.inspect

  if Array === attributes
    return attributes.map { |item_attributes| api_attributes(item_attributes) }
  end

  attributes["Schedules"].sort_by!{ |s| s["Kind"] } if attributes["Schedules"]

  # codes = attributes["Codes"]
  # if Array === codes
  #   attributes["Codes"] = Hash[codes.map { |code| [code["Kind"], code["Value"]] }]
  # end

  attributes
end

def has_attributes(response_array, attributes)
  parsed_attributes = model_attributes(attributes)

  code_space = parsed_attributes["Codes"].keys.first
  code_value = parsed_attributes["Codes"][code_space]

  found_value = response_array.find{|a| a["Codes"][code_space] == code_value}

  expect(found_value).not_to be_nil

  parsed_attributes.delete("Codes")

  parsed_attributes = parsed_attributes.reduce({}) do |attributes, (key, value)|
    case value
    when Float
      value = a_value_within(0.00001).of(value)
    when Array
      value = match_array(value)
    else
      value
    end

    attributes[key] = value
    attributes
  end

  expect(found_value).to include(parsed_attributes)
end

def gtfs_attributes(table)
  attributes = table.rows_hash
  attributes.each { |k, v| attributes[k] = eval("GTFS::Realtime::VehiclePosition::OccupancyStatus::#{v}") if k == "occupancy_status" }
end