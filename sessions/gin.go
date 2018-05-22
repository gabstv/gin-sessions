package sessions

import (
	"time"

	"github.com/gin-gonic/gin"
)

// MidCfg stores the configuration for the Middleware
type MidCfg struct {
	CookieVarName   string
	CookiePath      string
	CookieDomain    string
	CookieSecure    bool
	CookieHTTPOnly  bool
	NewSessionFunc  func(id string, data map[string]interface{}, expires time.Time) Session
	Logger          ErrorLogger
	DefaultDuration time.Duration
}

// Middleware is the session middleware
// To obtain the session:
// // session := c.MustGet("Session").(sessions.Session)
func Middleware(engine Engine, cfg *MidCfg) func(c *gin.Context) {
	if cfg == nil {
		cfg = &MidCfg{
			CookieVarName:   "_session_",
			NewSessionFunc:  Default,
			DefaultDuration: time.Minute * 30,
			CookieHTTPOnly:  true,
		}
	}
	return func(c *gin.Context) {
		sid := ""
		if ck, err := c.Cookie(cfg.CookieVarName); err == nil {
			sid = ck
		} else {
			if cfg.Logger != nil {
				cfg.Logger.Errorln("cookie error", cfg.CookieVarName, err.Error())
			}
		}
		if sid == "" || !engine.Exists(sid) {
			sid = NewID()
			sess := cfg.NewSessionFunc(sid, nil, time.Now().Add(cfg.DefaultDuration))
			c.Set("Session", sess)
			c.SetCookie(cfg.CookieVarName, sid, int(cfg.DefaultDuration/time.Second), cfg.CookiePath, cfg.CookieDomain, cfg.CookieSecure, cfg.CookieHTTPOnly)
		} else {
			sess, err := engine.Load(sid)
			if err != nil {
				if cfg.Logger != nil {
					cfg.Logger.Errorln("load session error", sid, err.Error())
				}
				sid = NewID()
				sess = cfg.NewSessionFunc(sid, nil, time.Now().Add(cfg.DefaultDuration))
			}
			c.Set("Session", sess)
			c.SetCookie(cfg.CookieVarName, sid, int(cfg.DefaultDuration/time.Second), cfg.CookiePath, cfg.CookieDomain, cfg.CookieSecure, cfg.CookieHTTPOnly)
		}
		c.Next()
		if isess, ok := c.Get("Session"); ok {
			if sess, ok2 := isess.(Session); ok2 {
				engine.Save(sess)
			}
		}
	}
}
