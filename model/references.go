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
