package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Model interface {
	Query(query string) ([]interface{}, error)
}

type Options struct {
	JSONRequestMaxBytes int64
}

type Router struct {
	routes  http.Handler
	options Options
	model   Model
}

func New(model Model) *Router {
	router := httprouter.New()
	rt := &Router{
		routes: router,
		options: Options{
			JSONRequestMaxBytes: 1_048_576,
		},
		model: model,
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
