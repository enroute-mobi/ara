def model_attributes(table)
  attributes = table.rows_hash.dup

  attributes.each do |key, value|
    case value
    when /\A\d+\Z/
      # Convert integer
      attributes[key] = value.to_i
    end
  end

  if attributes["ObjectIds"]
    attributes["ObjectIds"] = JSON.parse("{#{attributes["ObjectIds"]}}")
  end
  attributes
end
