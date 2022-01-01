package router

import (
	"errors"
	"fmt"
	"net/http"
)

func (rt *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	rt.notFoundResponse(w, r)
}

func (rt *Router) panicHandler(w http.ResponseWriter, r *http.Request, rcv interface{}) {
	err := errors.New(fmt.Sprintf("%v", rcv))
	rt.internalServerErrorResponse(w, r, err)
}
