package multiplexer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Multiplexer is a container of groups. This enables that same routes can have the same handler
// assigned like:
//     g1 := r.Group("/api")
//     g2 := r.Group("/:namespace/api")
//     m := multiplexer.New(g1, g2)
//     m.GET("/hello", helloHandler)
type Multiplexer interface {
	gin.IRoutes
	Group(string, ...gin.HandlerFunc) Multiplexer
}

type multiplexer struct {
	groups []*gin.RouterGroup
}

// New creates a new multiplex of
func New(groups ...*gin.RouterGroup) Multiplexer {
	grp := make([]*gin.RouterGroup, 0, len(groups))
	for _, v := range groups {
		grp = append(grp, v)
	}
	return &multiplexer{grp}
}

// Use adds middleware to all the groups
func (g *multiplexer) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.Use(middleware...)
	}
	return g
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in github.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (g *multiplexer) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.Handle(httpMethod, relativePath, handlers...)
	}
	return g
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *multiplexer) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.Any(relativePath, handlers...)
	}
	return g
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (g *multiplexer) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.GET(relativePath, handlers...)
	}
	return g
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (g *multiplexer) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.POST(relativePath, handlers...)
	}
	return g
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (g *multiplexer) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.DELETE(relativePath, handlers...)
	}
	return g
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (g *multiplexer) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.PATCH(relativePath, handlers...)
	}
	return g
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (g *multiplexer) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.PUT(relativePath, handlers...)
	}
	return g
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (g *multiplexer) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.OPTIONS(relativePath, handlers...)
	}
	return g
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (g *multiplexer) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	for _, v := range g.groups {
		v.HEAD(relativePath, handlers...)
	}
	return g
}

// StaticFile registers a single route in order to server a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (g *multiplexer) StaticFile(relativePath, filepath string) gin.IRoutes {
	for _, v := range g.groups {
		v.StaticFile(relativePath, filepath)
	}
	return g
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (g *multiplexer) Static(relativePath, root string) gin.IRoutes {
	for _, v := range g.groups {
		v.Static(relativePath, root)
	}
	return g
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (g *multiplexer) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	for _, v := range g.groups {
		v.StaticFS(relativePath, fs)
	}
	return g
}

// Group creates a new multiplex (of router groups). You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (g *multiplexer) Group(relativePath string, handlers ...gin.HandlerFunc) Multiplexer {
	newgroups := make([]*gin.RouterGroup, len(g.groups))
	for k, v := range g.groups {
		newgroups[k] = v.Group(relativePath, handlers...)
	}
	return &multiplexer{newgroups}
}
