package ledisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gabstv/gin-sessions/sessions"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

type driver int

var (
	newSessionM    sync.Mutex
	newSessionFunc func(id string, data map[string]interface{}, expires time.Time) sessions.Session
)

// SetNewSessionFunc changes the function that builds a new session
func SetNewSessionFunc(f func(id string, data map[string]interface{}, expires time.Time) sessions.Session) {
	newSessionM.Lock()
	defer newSessionM.Unlock()
	newSessionFunc = f
}

func (*driver) Start(ctx context.Context, args ...interface{}) (sessions.Engine, error) {

	newSessionM.Lock()
	if newSessionFunc == nil {
		newSessionFunc = sessions.Default
	}
	newSessionM.Unlock()

	if len(args) < 1 {
		return nil, fmt.Errorf("cannot initialize ledisdb session driver without args")
	}
	if ldb, ok := args[0].(*ledis.DB); ok {
		return &engine{
			l: ldb,
		}, nil
	}
	if ldcfg, ok := args[0].(*lediscfg.Config); ok {
		l, err := ledis.Open(ldcfg)
		if err != nil {
			return nil, err
		}
		db, err := l.Select(0)
		if err != nil {
			return nil, err
		}
		return &engine{
			l: db,
		}, nil
	}

	return nil, fmt.Errorf("cannot initialize ledisdb session driver (invalid args)")
}

type engine struct {
	l *ledis.DB
}

func (e *engine) Exists(id string) bool {
	if n, _ := e.l.Exists(K("session_%s", id)); n == 1 {
		return true
	}
	return false
}

func (e *engine) Load(id string) (sessions.Session, error) {
	b, err := e.l.Get(K("session_%s", id))
	if err != nil {
		return nil, err
	}
	mmm := make(map[string]interface{})
	err = json.Unmarshal(b[8:], &mmm)
	if err != nil {
		return nil, err
	}
	expu := ImportInt64(b[:8])
	newSessionM.Lock()
	sessgtr := newSessionFunc(id, mmm, time.Unix(expu, 0))
	newSessionM.Unlock()
	//
	return sessgtr, nil
}

func (e *engine) Save(s sessions.Session) error {
	mmm := s.GetAll()
	bb, err := json.Marshal(mmm)
	if err != nil {
		return err
	}
	ex := s.Expires()
	unix := ex.Unix()
	bb2 := make([]byte, len(bb)+8)
	copy(bb2, ExportInt64(unix))
	copy(bb2[8:], bb)

	err = e.l.Set(K("session_%s", s.ID()), bb2)
	if err != nil {
		return err
	}
	// get expiry
	sec := int64(ex.Sub(time.Now()).Seconds()) + 1
	e.l.Expire(K("session_%s", s.ID), sec)
	return nil
}

func (e *engine) Count() int {
	var iter []byte
	ok := true
	n := 0
	for ok {
		res, err := e.l.Scan(ledis.KV, iter, 10, false, "session_.*")
		if err != nil {
			ok = false
			continue
		}
		if len(res) < 1 {
			ok = false
			continue
		}
		iter = res[len(res)-1]
		n += len(res)
	}
	return n
}

func init() {
	var d driver = 1
	sessions.Register("ledisdb", &d)
}
