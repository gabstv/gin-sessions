package sessions

import (
	"context"
	"fmt"
	"sync"
)

var (
	drivers  = make(map[string]Driver)
	driversM sync.Mutex
)

// Engine represents a session engine.
type Engine interface {
	Exists(id string) bool
	Load(id string) (Session, error)
	Save(s Session) error
	Count() int
}

// Driver represents a session engine driver.
type Driver interface {
	Start(ctx context.Context, args ...interface{}) (Engine, error)
}

// Register registers a session driver.
func Register(name string, driver Driver) error {
	driversM.Lock()
	defer driversM.Unlock()
	if _, exists := drivers[name]; exists {
		return fmt.Errorf("session driver %v already registered", name)
	}
	drivers[name] = driver
	return nil
}

// Connect connects to a registered session driver.
func Connect(ctx context.Context, name string, args ...interface{}) (Engine, error) {
	driversM.Lock()
	defer driversM.Unlock()
	if d, ok := drivers[name]; ok {
		return d.Start(ctx, args...)
	}
	return nil, fmt.Errorf("driver %s not found", name)
}
