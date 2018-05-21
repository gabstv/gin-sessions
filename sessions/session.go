package sessions

import (
	"sync"
	"time"
)

// Session is what is saved on the server. The browser stores the ID
// inside a cookie like sessid.
type Session interface {
	ID() string
	Expires() time.Time
	SetExpires(t time.Time)
	Get(key string) interface{}
	Set(key string, value interface{})
	GetAll() map[string]interface{}
}

type defaultsess struct {
	sync.RWMutex
	m  map[string]interface{}
	ex time.Time
	id string
}

// Default returns a basic Session implementation
func Default(id string, data map[string]interface{}, expires time.Time) Session {
	return &defaultsess{
		m:  data,
		ex: expires,
		id: id,
	}
}

func (s *defaultsess) ID() string {
	s.RLock()
	defer s.RUnlock()
	return s.id
}

func (s *defaultsess) Expires() time.Time {
	s.RLock()
	defer s.RUnlock()
	return s.ex
}

func (s *defaultsess) SetExpires(t time.Time) {
	s.Lock()
	defer s.Unlock()
	s.ex = t
}

func (s *defaultsess) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	if s.m == nil {
		return nil
	}
	return s.m[key]
}

func (s *defaultsess) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	if s.m == nil {
		s.m = make(map[string]interface{})
	}
	s.m[key] = value
}

func (s *defaultsess) GetAll() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	if s.m == nil {
		return nil
	}
	m2 := make(map[string]interface{})
	for k, v := range s.m {
		m2[k] = v
	}
	return m2
}
