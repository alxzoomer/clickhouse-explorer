package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Options struct {
	JSONRequestMaxBytes int64
}

type Router struct {
	routes  http.Handler
	options Options
}

func New() *Router {
	router := httprouter.New()
	rt := &Router{
		routes: router,
		options: Options{
			JSONRequestMaxBytes: 1_048_576,
		},
	}
	router.NotFound = http.HandlerFunc(rt.notFoundHandler)
	router.PanicHandler = rt.panicHandler
	router.GET("/", rt.indexHandler)
	router.POST("/api/v1/query", rt.queryHandler)
	return rt
}

func (rt *Router) Handler() http.Handler {
	return rt.routes
}
