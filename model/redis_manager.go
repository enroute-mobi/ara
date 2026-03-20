package model

import (
	"encoding/json"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

/*
To be used as an embedded struct.
P is a model struct, and T its pointer

	Example:

	type LineManager struct {
		redisManager[LineId, Line, *Line]
	}
*/
type redisManager[Id ~string, P any, T RedisModelInstance[P]] struct {
	uuid.UUIDConsumer

	new       func(Model) T
	modelName string
	model     Model
	client    redisclient.Client
}

func (m *redisManager[Id, P, T]) New() T {
	return m.new(m.model)
}

func (m *redisManager[Id, P, T]) Find(id Id) (T, bool) {
	val, err := m.client.Get(m.modelName, string(id))
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v: %v", m.modelName, id, err)
		return nil, false
	}

	t := m.new(m.model)
	err = json.Unmarshal([]byte(val), t)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v: %v", m.modelName, id, err)
		return nil, false
	}
	return t, true
}

func (m *redisManager[Id, P, T]) FindAttribute(id Id, attr string) (string, bool) {
	b := strings.Builder{}
	b.WriteString("$.")
	b.WriteString(attr)

	val, err := m.client.GetPath(m.modelName, string(id), b.String())
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attribute %v: %v", m.modelName, id, attr, err)
		return "", false
	}

	t := make(map[string]string)
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attribute %v: %v", m.modelName, id, attr, err)
		return "", false
	}
	return t[b.String()], true
}

func (m *redisManager[Id, P, T]) FindAttributes(id Id, attrs ...string) (map[string]string, bool) {
	if len(attrs) == 0 {
		return nil, false
	}

	b := strings.Builder{}
	for i := range attrs {
		b.WriteString("$.")
		b.WriteString(attrs[i])
		b.WriteString(" ")
	}

	val, err := m.client.GetPath(m.modelName, string(id), b.String())
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attributes %v: %v", m.modelName, id, attrs, err)
		return nil, false
	}

	t := make(map[string]string)
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Log.Debugf("Error while finding %v %v attributes %v: %v", m.modelName, id, attrs, err)
		return nil, false
	}

	r := make(map[string]string)
	for k, v := range t {
		r[k[2:]] = v // Remove first two characters of the key "$."
	}
	return r, true
}

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

func (m *redisManager[Id, P, T]) FindBy(indexName, identifier string) (T, bool) {
	docs, err := m.client.FindBy(m.modelName, indexName, identifier)
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

func (m *redisManager[Id, P, T]) FindAllBy(indexName, identifier string) []T {
	docs, err := m.client.FindBy(m.modelName, indexName, identifier)
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

func (m *redisManager[Id, P, T]) Save(t T) bool {
	err := m.client.Set(t)
	return err == nil
}

func (m *redisManager[Id, P, T]) Delete(t T) bool {
	err := m.client.Delete(t.ModelId())
	return err == nil
}

func (m *redisManager[Id, P, T]) SetModel(model Model) {
	m.model = model
}

/*
   For managers which handles code
   If we need to have more embedded structs like that, we should use another pattern:

   type codeHandler[Id ~string, P any, T RedisModelInstance[P]] struct{
     rm *RedisManager[Id, P, T]
   }

	type LineManager struct {
		redisManager[LineId, Line, *Line]
		codeHandler
	}

    manager := LineManager{}
	manager.codeHandler.rm = &(manager.redisManager)
*/

type redisCodeHandlerManager[Id ~string, P any, T RedisModelInstance[P]] struct {
	redisManager[Id, P, T]
}

func (m *redisCodeHandlerManager[Id, P, T]) FindByCode(c Code) (T, bool) {
	docs, err := m.client.FindByCode(m.modelName, c.CodeSpace(), c.Value())
	if err != nil {
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

func (m *redisCodeHandlerManager[Id, P, T]) CodeExists(c Code) bool {
	docs, err := m.client.FindByCode(m.modelName, c.CodeSpace(), c.Value())
	if err != nil {
		logger.Log.Debugf("Error While finding %v by code: %v", m.modelName, err)
		return false
	}

	return len(docs) == 1
}

func (m *redisCodeHandlerManager[Id, P, T]) FindCode(id Id, codespace string) (Code, bool) {
	b := strings.Builder{}
	b.WriteString("$.codes.")
	b.WriteString(codespace)
	b.WriteString(".value")

	val, err := m.client.GetPath(m.modelName, string(id), b.String())
	if err != nil {
		logger.Log.Debugf("Error while finding %v code %v attribute %v: %v", m.modelName, id, codespace, err)
		return Code{}, false
	}

	t := make(map[string]string)
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		logger.Log.Debugf("Error while finding %v code %v attribute %v: %v", m.modelName, id, codespace, err)
		return Code{}, false
	}
	return NewCode(codespace, t[b.String()]), true
}
