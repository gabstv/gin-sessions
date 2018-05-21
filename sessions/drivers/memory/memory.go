package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gabstv/gin-sessions/sessions"
)

type driver int

func (*driver) Start(ctx context.Context, args ...interface{}) (sessions.Engine, error) {
	entries := make(map[string]sessions.Session)
	eng := &engine{
		entries: entries,
	}
	go func(ctx context.Context, e *engine) {
		for ctx.Err() == nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 30):
			}
			time.Sleep(time.Second * 30)
			e.Lock()
			n := time.Now()
			delist := make([]string, 0)
			for k, v := range e.entries {
				if v.Expires().Before(n) {
					delist = append(delist, k)
				}
			}
			for _, v := range delist {
				delete(e.entries, v)
			}
			e.Unlock()
		}
	}(ctx, eng)
	return eng, nil
}

type engine struct {
	sync.Mutex
	entries map[string]sessions.Session
}

func (e *engine) Exists(id string) bool {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.entries[id]; ok {
		return true
	}
	return false
}

func (e *engine) Load(id string) (sessions.Session, error) {
	e.Lock()
	defer e.Unlock()
	if s, ok := e.entries[id]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("not found")
}

func (e *engine) Save(s sessions.Session) error {
	e.Lock()
	defer e.Unlock()
	e.entries[s.ID()] = s
	return nil
}

func (e *engine) Count() int {
	e.Lock()
	defer e.Unlock()
	return len(e.entries)
}

func init() {
	var d driver = 1
	sessions.Register("memory", &d)
}
