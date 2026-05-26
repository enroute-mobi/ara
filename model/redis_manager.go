package model

import (
	"encoding/json"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

// redisManager defines a generic manager using Redis.
// It handles all the basic commands like New, Find, FindBy, etc...
//
// It needs to be used as an embedded struct.
// Id is an Id type like LineId, P a model struct, and T its pointer
// T needs to satisfy the RedisModelInstance interface, meaning it
// needs to have a SetId(Id) and a ModelId() method.
//
// Example:
//
//	type LineManager struct {
//		redisManager[LineId, Line, *Line]
//	}
//
// We need to provide it with:
// new: a method to create an object of type T
// modelName: a string containing the modelName to use with the redis.Client
type redisManager[Id ~string, P any, T RedisModelInstance[Id, P]] struct {
	uuid.UUIDConsumer

	new       func(Model) T
	modelName string
	model     Model
	client    redisclient.Client
}

// New returns a new record properly initialized
func (m *redisManager[Id, P, T]) New() T {
	return m.new(m.model)
}

// Find returns a record by its Id and true if we found it, otherwise nil and false
func (m *redisManager[Id, P, T]) Find(id Id) (T, bool) {
	val, err := m.client.Get(m.modelName, string(id))
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v: %v", m.modelName, id, err)
		return nil, false
	}

	ts := []T{m.new(m.model)}
	err = json.Unmarshal([]byte(val), &ts)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v: %v", m.modelName, id, err)
		return nil, false
	}
	return ts[0], true
}

// FindAttribute finds a record by its Id and returns only the requested attribute.
// It returns true if we found the record
func (m *redisManager[Id, P, T]) FindAttribute(id Id, attr string) (string, bool) {
	b := strings.Builder{}
	b.WriteString("$.")
	b.WriteString(attr)

	val, err := m.client.Get(m.modelName, string(id), b.String())
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attribute %v: %v", m.modelName, id, attr, err)
		return "", false
	}

	t := []string{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil || len(t) != 1 {
		logger.Log.Debugf("Error while finding %v %v attribute %v: %v", m.modelName, id, attr, err)
		return "", false
	}
	return t[0], true
}

// FindAttribute finds a record by its Id and returns all the requested attributes
// in a map[string]any. It returns true if we found the record
func (m *redisManager[Id, P, T]) FindAttributes(id Id, attrs ...string) (map[string]any, bool) {
	// If there's no attributes, the method can't properly work, but it's an internal method
	// and checking this costs time.
	// if len(attrs) == 0 {
	// 	return nil, false
	// }

	b := strings.Builder{}
	for i := range attrs {
		b.Reset()
		b.WriteString("$.")
		b.WriteString(attrs[i])
		attrs[i] = b.String()
	}

	val, err := m.client.Get(m.modelName, string(id), attrs...)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attributes %v: %v", m.modelName, id, attrs, err)
		return nil, false
	}

	t := make(map[string][]any)
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attributes %v: %v", m.modelName, id, attrs, err)
		return nil, false
	}

	r := make(map[string]any)
	for k, v := range t {
		if len(v) != 1 {
			continue
		}
		r[k[2:]] = v[0] // Remove first two characters of the key: "$."
	}
	return r, true
}

// FindAll returns all records
func (m *redisManager[Id, P, T]) FindAll() []T {
	docs, err := m.client.FindAll(m.modelName)
	if err != nil {
		logger.Log.Debugf("Error While finding all %v: %v", m.modelName, err)
		return nil
	}

	ts := []T{}
	for i := range docs {
		t := m.new(m.model)
		err = json.Unmarshal([]byte(docs[i].Fields["$"]), t)
		if err != nil {
			logger.Log.Debugf("Error While finding all %v: %v", m.modelName, err)
			continue
		}
		ts = append(ts, t)
	}

	return ts
}

// FindBy returns one record, finding it by an attribute we have an index on
func (m *redisManager[Id, P, T]) FindBy(indexName, value string) (T, bool) {
	docs, err := m.client.FindBy(m.modelName, indexName, value)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}
	if l := len(docs); l != 1 {
		logger.Log.Debugf("Error While finding %v by %v: returned %v instances: %v", m.modelName, indexName, l, docs)
		return nil, false
	}

	t := m.new(m.model)
	err = json.Unmarshal([]byte(docs[0].Fields["$"]), t)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}

	return t, true
}

// FindAttributesBy finds one record by an attribute we have an index on,
// and returns a map[string]string (map[attribute]value)
func (m *redisManager[Id, P, T]) FindAttributesBy(indexName, value string, attrs ...string) (map[string]string, bool) {
	if len(attrs) == 0 {
		return nil, false
	}

	docs, err := m.client.FindBy(m.modelName, indexName, value, attrs...)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}
	if l := len(docs); l != 1 {
		logger.Log.Debugf("Error While finding %v by %v: returned %v instances: %v", m.modelName, indexName, l, docs)
		return nil, false
	}

	return docs[0].Fields, true
}

// FindAllBy returns an array of records, finding them by an attribute we have an index on
func (m *redisManager[Id, P, T]) FindAllBy(indexName, value string) []T {
	docs, err := m.client.FindBy(m.modelName, indexName, value)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil
	}

	ts := []T{}
	for i := range docs {
		t := m.new(m.model)
		err = json.Unmarshal([]byte(docs[i].Fields["$"]), t)
		if err != nil {
			logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
			continue
		}
		ts = append(ts, t)
	}

	return ts
}

// FindAllAttributesBy finds multiple records by an attribute we have an index on,
// and returnsan array of map[string]string (map[attribute]value)
func (m *redisManager[Id, P, T]) FindAllAttributesBy(indexName, value string, attrs ...string) []map[string]string {
	docs, err := m.client.FindBy(m.modelName, indexName, value, attrs...)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil
	}

	ts := []map[string]string{}
	for i := range docs {
		ts = append(ts, docs[i].Fields)
	}

	return ts
}

func (m *redisManager[Id, P, T]) Save(t T) bool {
	if t.ModelId() == "" {
		t.SetId(Id(m.NewUUID()))
	}
	err := m.client.Set(m.modelName, t)
	return err == nil
}

func (m *redisManager[Id, P, T]) Delete(t T) bool {
	err := m.client.Delete(m.modelName, t.ModelId())
	return err == nil
}

func (m *redisManager[Id, P, T]) SetModel(model Model) {
	m.model = model
}

// redisCodeHandlerManager adds code handling to redisManager
// If we need to have more embedded structs like that, we should use another pattern:
//
//	type codeHandler[Id ~string, P any, T RedisModelInstance[Id, P]] struct{
//		rm *RedisManager[Id, P, T]
//	}
//
//	type LineManager struct {
//		redisManager[LineId, Line, *Line]
//		codeHandler
//	}
//
// manager := LineManager{}
// manager.codeHandler.rm = &(manager.redisManager)
type redisCodeHandlerManager[Id ~string, P any, T RedisModelInstance[Id, P]] struct {
	redisManager[Id, P, T]
}

// FindByCode returns one record, finding it by an code.
// We need to have an index on the code space
func (m *redisCodeHandlerManager[Id, P, T]) FindByCode(c Code) (T, bool) {
	docs, err := m.client.FindByCode(m.modelName, c.CodeSpace(), c.Value())
	if err != nil || len(docs) == 0 {
		logger.Log.Debugf("Error While finding %v by code: %v", m.modelName, err)
		return nil, false
	}

	t := m.new(m.model)
	err = json.Unmarshal([]byte(docs[0].Fields["$"]), t)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by code: %v", m.modelName, err)
		return nil, false
	}

	return t, true

}

// CodeExists returns true if we can find a record withe the given code
func (m *redisCodeHandlerManager[Id, P, T]) CodeExists(c Code) bool {
	docs, err := m.client.FindByCode(m.modelName, c.CodeSpace(), c.Value())
	if err != nil {
		logger.Log.Debugf("Error While finding %v by code: %v", m.modelName, err)
		return false
	}

	return len(docs) == 1
}

// FindCode finds a record by its Id, and returns a code if it exists for the given code space
func (m *redisCodeHandlerManager[Id, P, T]) FindCode(id Id, codespace string) (Code, bool) {
	b := strings.Builder{}
	b.WriteString("$.Codes.")
	b.WriteString(codespace)

	val, err := m.client.Get(m.modelName, string(id), b.String())
	if err != nil {
		logger.Log.Debugf("Error while finding %v code %v attribute %v: %v", m.modelName, id, codespace, err)
		return Code{}, false
	}

	t := []string{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil || len(t) != 1 {
		logger.Log.Debugf("Error while finding %v code %v attribute %v: %v", m.modelName, id, codespace, err)
		return Code{}, false
	}
	return NewCode(codespace, t[0]), true
}
