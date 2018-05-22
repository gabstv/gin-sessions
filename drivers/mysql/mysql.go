package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gabstv/gin-sessions/sessions"
	"github.com/gabstv/sqltypes"
)

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

type driver int

func runCleanup(ctx context.Context, db *sql.DB, tblprefix string) {
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute * 15):
		}
		nnow := time.Now().Unix()
		//
		db.Exec(fmt.Sprintf("DELETE FROM '%sentries' WHERE expires < ?", tblprefix), nnow)
	}
}

func (*driver) Start(ctx context.Context, args ...interface{}) (sessions.Engine, error) {

	tblprefix := "_sessgin_"

	newSessionM.Lock()
	if newSessionFunc == nil {
		newSessionFunc = sessions.Default
	}
	newSessionM.Unlock()

	if len(args) < 1 {
		return nil, fmt.Errorf("cannot initialize mysql session driver without args")
	}
	if db, ok := args[0].(*sql.DB); ok {
		if len(args) > 1 {
			if str, ok2 := args[1].(string); ok2 {
				if str == "cleanup=true" {
					//TODO: config for tblprefix
					go runCleanup(ctx, db, tblprefix)
				}
			}
		}
		return &engine{
			db:     db,
			prefix: tblprefix,
		}, nil
	}
	//TODO: support more entry methods

	return nil, fmt.Errorf("cannot initialize ledisdb session driver (invalid args)")
}

type engine struct {
	db     *sql.DB
	prefix string
}

func (e *engine) Exists(id string) bool {
	var n int
	e.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM '%sentries' WHERE id=?", e.prefix), id).Scan(&n)
	return n > 0
}

func (e *engine) Load(id string) (sessions.Session, error) {
	var data sqltypes.NullString
	var expires int64
	if err := e.db.QueryRow(fmt.Sprintf("SELECT data, expires FROM '%sentries' WHERE id=?", e.prefix), id).
		Scan(&data, &expires); err != nil {
		return nil, err
	}
	mm := make(map[string]interface{})
	if data != "" {
		if err := json.Unmarshal([]byte(data.String()), &mm); err != nil {
			return nil, err
		}
	}
	newSessionM.Lock()
	s := newSessionFunc(id, mm, time.Unix(expires, 0))
	newSessionM.Unlock()
	return s, nil
}

func (e *engine) Save(s sessions.Session) error {
	e64 := s.Expires().Unix()
	d := s.GetAll()
	var ds sqltypes.NullString
	if d != nil {
		b, err := json.Marshal(&d)
		if err != nil {
			return err
		}
		ds = sqltypes.NullString(string(b))
	}
	_, err := e.db.Exec(fmt.Sprintf("INSERT INTO '%sentries' (id, data, expires) VALUES (?,?,?) ON DUPLICATE KEY UPDATE data=?, expires=?", e.prefix), s.ID(), ds, e64, ds, e64)
	return err
}

func (e *engine) Count() int {
	var n int
	e.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM '%sentries'", e.prefix)).Scan(&n)
	return n
}

func init() {
	var d driver = 1
	sessions.Register("mysql", &d)
}
