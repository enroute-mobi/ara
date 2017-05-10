package model

type References map[string]Reference

func NewReferences() References {
	return make(References)
}

func (references References) Set(key string, value Reference) {
	emptyRef := Reference{} // Compile error...
	if value == emptyRef {
		return
	}
	if value.ObjectId.Kind() == "" || value.ObjectId.Value() == "" {
		return
	}
	references[key] = value
}

func (references References) SetObjectId(key string, obj ObjectID, id string) {
	if obj.Kind() == "" || obj.Value() == "" {
		return
	}
	references[key] = Reference{ObjectId: &obj, Id: id}
}

func (references References) IsEmpty() bool {
	return len(references) == 0
}

func (references References) Copy() References {
	copy := NewReferences()

	for key, value := range references {
		var obj ObjectID
		if value.ObjectId != nil {
			obj = NewObjectID(value.ObjectId.Value(), value.ObjectId.Kind())
		}
		copy[key] = Reference{
			Id:       value.Id,
			Type:     value.Type,
			ObjectId: &obj,
		}
	}
	return copy
}
