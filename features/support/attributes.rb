def model_attributes(table)
  attributes = table.rows_hash
  if attributes["ObjectIds"]
    attributes["ObjectIds"] = JSON.parse("{#{attributes["ObjectIds"]}}")
  end
  attributes
end
