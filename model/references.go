package model

import "sync"

type References struct {
	ref   map[string]Reference
	mutex *sync.RWMutex
}

func NewReferences() References {
	return References{
		ref:   make(map[string]Reference),
		mutex: &sync.RWMutex{},
	}
}

func (references *References) GetReferences() map[string]Reference {
	return references.ref
}

func (references *References) GetSiriReferences() map[string]string {
	sref := make(map[string]string)

	references.mutex.RLock()
	for k, v := range references.ref {
		sref[k] = v.ObjectId.Value()
	}
	references.mutex.RUnlock()
	return sref
}

func (references *References) SetReference(key string, r Reference) {
	references.mutex.Lock()
	defer references.mutex.Unlock()

	references.ref[key] = r
}

func (references *References) SetReferences(r map[string]Reference) {
	references.ref = r
}

func (references *References) Get(key string) (Reference, bool) {
	references.mutex.RLock()
	defer references.mutex.RUnlock()

	ref, ok := references.ref[key]
	return ref, ok
}

func (references *References) Set(key string, value Reference) {
	references.mutex.Lock()
	defer references.mutex.Unlock()

	emptyRef := Reference{} // Compile error...
	if value == emptyRef {
		return
	}
	if value.ObjectId.Kind() == "" || value.ObjectId.Value() == "" {
		return
	}
	references.ref[key] = value
}

func (references *References) SetObjectId(key string, obj ObjectID) {
	references.mutex.Lock()
	defer references.mutex.Unlock()

	if obj.Kind() == "" || obj.Value() == "" {
		return
	}
	references.ref[key] = Reference{ObjectId: &obj}
}

func (references *References) IsEmpty() bool {
	return len(references.ref) == 0
}

func (references *References) Copy() References {
	references.mutex.Lock()
	defer references.mutex.Unlock()

	newReferences := NewReferences()

	for key, value := range references.ref {
		newReferences.ref[key] = Reference{
			Type: value.Type,
		}
		if value.ObjectId != nil {
			objectid := NewObjectID(value.ObjectId.Kind(), value.ObjectId.Value())
			tmp := newReferences.ref[key]
			tmp.ObjectId = &objectid
			newReferences.ref[key] = tmp
		}
	}
	return newReferences
}
