package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	routes http.Handler
}

func New() *Router {
	router := httprouter.New()
	rt := &Router{
		routes: router,
	}
	router.NotFound = http.HandlerFunc(rt.notFoundHandler)
	router.PanicHandler = rt.panicHandler
	router.GET("/", rt.indexHandler)
	router.GET("/api/v1/query", rt.queryHandler)
	return rt
}

func (rt *Router) Handler() http.Handler {
	return rt.routes
}
