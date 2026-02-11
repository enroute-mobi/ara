package settings

import (
	"encoding/json"
	"maps"
	"sync"
)

type Settings struct {
	m *sync.RWMutex

	r func()

	s map[string]string
}

func NewSettings() Settings {
	return Settings{
		m: &sync.RWMutex{},
		r: func() {},
		s: make(map[string]string),
	}
}

func (s *Settings) SettingsLen() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.s)
}

func (s *Settings) Setting(key string) string {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.s[key]
}

// Should only be used in tests
func (s *Settings) SetSetting(k, v string) {
	s.m.Lock()
	s.s[k] = v
	s.r()
	s.m.Unlock()
}

func (s *Settings) SettingsDefinition() (m map[string]string) {
	m = make(map[string]string)
	s.m.RLock()
	maps.Copy(m, s.s)
	s.m.RUnlock()
	return
}

func (s *Settings) SetSettingsDefinition(m map[string]string) {
	if m == nil {
		return
	}
	s.m.Lock()
	s.s = make(map[string]string)
	maps.Copy(s.s, m)
	s.r()
	s.m.Unlock()
}

func (s *Settings) ToJson() ([]byte, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return json.Marshal(s.s)
}
