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
		sref[k] = v.Code.Value()
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
	if value.Code.CodeSpace() == "" || value.Code.Value() == "" {
		return
	}
	references.ref[key] = value
}

func (references *References) SetCode(key string, obj Code) {
	references.mutex.Lock()
	defer references.mutex.Unlock()

	if obj.CodeSpace() == "" || obj.Value() == "" {
		return
	}
	references.ref[key] = Reference{Code: &obj}
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
		if value.Code != nil {
			code := NewCode(value.Code.CodeSpace(), value.Code.Value())
			tmp := newReferences.ref[key]
			tmp.Code = &code
			newReferences.ref[key] = tmp
		}
	}
	return newReferences
}
