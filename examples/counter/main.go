package main

import (
	"context"
	"net/http"

	_ "github.com/gabstv/gin-sessions/drivers/memory"
	"github.com/gabstv/gin-sessions/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancelf := context.WithCancel(context.Background())
	defer cancelf()
	r := gin.Default()

	engine, err := sessions.Connect(ctx, "memory")
	if err != nil {
		panic(err)
	}

	g0 := r.Group("/without_session")
	g0.GET("/count", func(c *gin.Context) {
		// session := c.MustGet("Session").(sessions.Session) // this would panic
		if isess, ok := c.Get("Session"); ok {
			if session, ok2 := isess.(sessions.Session); ok2 {
				cn, _ := session.Get("count").(int)
				session.Set("count", cn+1)
				c.String(http.StatusOK, "Count: %d", cn)
				return
			}
		}
		c.String(http.StatusOK, "Count: ?")
	})

	g1 := r.Group("/with_session")
	g1.Use(sessions.Middleware(engine, nil))
	g1.GET("/count", func(c *gin.Context) {
		session := c.MustGet("Session").(sessions.Session)
		cn, _ := session.Get("count").(int)
		session.Set("count", cn+1)
		c.String(http.StatusOK, "Count: %d", cn)
	})

	r.Run(":7766")
}
