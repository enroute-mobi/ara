package model

import (
	"encoding/json"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	code = "code"
)

type RedisManager[Id string, P any, T RedisModelInstance[P]] struct {
	uuid.UUIDConsumer

	modelName string
	model     Model
	client    *redisclient.Client
}

func (m *RedisManager[Id, P, T]) New() T {
	return new(P)
}

func (m *RedisManager[Id, P, T]) Find(id Id) (T, error) {
	val, err := m.client.Get(m.prefix(string(id)))
	if err != nil {
		return nil, err
	}

	t := new(P)
	err = json.Unmarshal([]byte(val), t)
	return t, err
}

func (m *RedisManager[Id, P, T]) FindAll() ([]T, error) {
	val, err := m.client.Get(m.modelName)
	if err != nil {
		return nil, err
	}

	t := make([]T, 300)
	err = json.Unmarshal([]byte(val), &t)
	return t, err
}

func (m *RedisManager[Id, P, T]) FindByCode(c Code) (T, bool) {
	return m.FindOneBy(code, c.String())
}

func (m *RedisManager[Id, P, T]) CodeExists(c Code) bool {
	docs, err := m.client.FindBy(m.modelName, code, c.String())
	if err != nil {
		logger.Log.Debugf("Error While finding %v by code: %v", m.modelName, err)
		return false
	}
	if l := len(docs); l == 1 {
		return true
	}
	return false
}

func (m *RedisManager[Id, P, T]) FindOneBy(indexName, identifier string) (T, bool) {
	docs, err := m.client.FindBy(m.modelName, indexName, identifier)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}
	if l := len(docs); l != 1 {
		logger.Log.Debugf("Error While finding %v by %v: returned %v instances: %v", m.modelName, indexName, l, docs)
		return nil, false
	}

	t := new(P)
	err = json.Unmarshal([]byte(docs[0].Fields["$"]), t)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}

	return t, true
}

func (m *RedisManager[Id, P, T]) FindBy(indexName, identifier string) ([]T, bool) {
	docs, err := m.client.FindBy(m.modelName, indexName, identifier)
	if err != nil {
		logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
		return nil, false
	}

	ts := []T{}
	for i := range docs {
		t := new(P)
		err = json.Unmarshal([]byte(docs[i].Fields["$"]), t)
		if err != nil {
			logger.Log.Debugf("Error While finding %v by %v: %v", m.modelName, indexName, err)
			continue
		}
		ts = append(ts, t)
	}

	return ts, true
}

func (m *RedisManager[Id, P, T]) Save(t T) bool {
	err := m.client.Save(t)
	if err != nil {
		return false
	}
	return true
}

func (m *RedisManager[Id, P, T]) Delete(t T) bool {
	err := m.client.Save(t)
	if err != nil {
		return false
	}
	return true
}

func (m *RedisManager[Id, P, T]) SetModel(model Model) {
	m.model = model
}

func (m *RedisManager[Id, P, T]) prefix(s string) string {
	b := strings.Builder{}
	b.Grow(50)
	b.WriteString(m.modelName)
	b.WriteRune(':')
	b.WriteString(s)
	return b.String()
}
