def model_attributes(table)
  attributes = table.rows_hash.dup

  attributes.dup.each do |key, value|
    case value
    when "nil"
      attributes[key] = nil
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
    # | Blocking[JourneyPlanner]  | true  |
    # | Blocking[RealTime]        | false |
    # into
    # "Blocking" => {"JourneyPlanner" => true, "Realtime" => false }
    if key =~ /Blocking\[([^\]]+)\]/
      name = $1
      attributes["Blocking"] ||= {}
      attributes["Blocking"][name] = value == "true" unless attributes["Blocking"].key?(name)
      attributes.delete key
    end

    # Transform
    # | KEY[A]  | value1 |
    # | KEY[B]  | value2 |
    # into
    # "KEY" => {"A" => "value1", "B" => "value2" }
    # With KEY either Codes, Attributes
    if key =~ /(Codes|Attributes)\[([^\]]+)\]/
      attr = Regexp.last_match(1)
      name = Regexp.last_match(2)

      attributes[attr] ||= {}
      attributes[attr][name] = value
      attributes.delete key
    end

    # Transform
    #  | KEY[A]     | value1 |
    #  | KEY[B]#FR  | value2 |
    # into
    # "KEY" => {"A"=>"value1", "B"=> { "FR" => "value2" } }
    #
    # With KEY either Summary, Description or Prompt
    # and B either DefaultValue, Translations
    if key =~ /(Summary|Description|Prompt)\[(DefaultValue|Translations)\](#(\S+))?/
      text_type = Regexp.last_match(1)
      name = Regexp.last_match(2)
      language = Regexp.last_match(4)
      case name
      when 'DefaultValue'
        attributes[text_type] ||= {}
        attributes[text_type][name] = value
      when 'Translations'
        attributes[text_type] ||= {}
        attributes[text_type][name] ||= {}
        attributes[text_type][name][language] = value
      end
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

      attributes[name][period_number][attribute.to_s] = if value == "nil"
                                                          nil
                                                        else
                                                          value
                                                        end

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

    if key =~ %r{^(Affects\[([^\]]+)\])(/(AffectedDestinations|AffectedSections|AffectedRoutes|LineIds)\[(\d+)])?(/((FirstStop|LastStop)|StopAreaId\z|RouteRef|StopAreaIds))?}
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
      when 'AllLines'
        attributes['Affects'] << { 'Type' => attribute }
      else
        model, id = attribute.split('=')
        attributes['Affects'].map do |affect|
          next if affect['Type'] != model && affect["#{model}Id"] != id

          affect[subaffect] ||= []
          if stop_type == 'StopAreaIds'
            affect[subaffect][index]['StopAreaIds'] ||= []
            affect[subaffect][index]['StopAreaIds'] << value
          elsif stop_type != ""
            affect[subaffect][index] ||= {}
            affect[subaffect][index][stop_type] = value
          else
            affect[subaffect] ||= []
            affect[subaffect][index] = value
          end
        end
      end
      attributes.delete key
    end

    # Transform
    # | ReferenceArray[0] | A, "B": "C" |
    # into
    # References => [ { "Type" => A, "Code" => { "B" => "C" } } ]
    # or
    # | ReferenceArray[0] | "B": "C" |
    # into
    # References => [ { "Code" => { "B" => "C" } } ]
    if key =~ /ReferenceArray\[(\d+)\]/
      name = $1
      attribute = $2

      attributes["References"] ||= []

      values = value.split(',')

      attributes["References"][$1.to_i] = {
        "Code" => JSON.parse("{ #{values[1]} }")
      }

      if values.size > 1
        attributes["References"][$1.to_i]["Type"] = values[0]
      end

      attributes.delete key
    end
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

def matcher_attributes(attributes, found_model)
  parsed_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  parsed_attributes = parsed_attributes.each_with_object({}) do |(key, value), attributes|
    case value
    when Float
      value = a_value_within(0.00001).of(value)
    when Array
      value = match_array(value)
    else
      value
    end

    attributes[key] = value
  end

  expect(found_model).to have_attributes(parsed_attributes)
end

def check_attributes(referential, model, attributes)
  parsed_attributes = model_attributes(attributes).transform_keys { |key| key.to_s.underscore.to_sym }

  code_space = parsed_attributes[:codes].keys.first
  value = parsed_attributes[:codes][code_space]

  found_model = referential.send(model).find("#{code_space}:#{value}")
  expect(found_model).not_to be_nil

  parsed_attributes.delete(:codes)

  parsed_attributes = parsed_attributes.each_with_object({}) do |(key, value), attributes|
    case value
    when Float
      value = a_value_within(0.00001).of(value)
    when Array
      value = match_array(value)
    else
      value
    end

    attributes[key] = value
  end

  expect(found_model).to have_attributes(parsed_attributes)
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
  attributes = table.rows_hash.dup

  attributes.dup.each do |key, value|
    if key =~ /direction_id/
      attributes[key] = value.to_i
    end
  end

  attributes.each { |k, v| attributes[k] = eval("GTFS::Realtime::VehiclePosition::OccupancyStatus::#{v}") if k == "occupancy_status" }

  attributes.dup.each do |key, value|
    if key =~ /(header_text_translation|description_text_translation)\[([^\]]+)\]/

      name = Regexp.last_match(1)
      lang = JSON.parse(Regexp.last_match(2))
      attributes[name] ||= {}
      attributes[name][lang] = value


      attributes.delete key
    end


  end

  attributes
end
