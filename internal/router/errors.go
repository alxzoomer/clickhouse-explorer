package router

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime/debug"
)

func (rt *Router) logError(err error, r *http.Request, status int) {
	log.Error().
		Err(err).
		Str("method", r.Method).
		Stringer("url", r.URL).
		Int("status", status).
		Bytes("stack", debug.Stack()).
		Msg("")
}

func (rt *Router) logWarn(r *http.Request, status int) {
	log.Warn().
		Str("method", r.Method).
		Stringer("url", r.URL).
		Int("status", status).
		Msg("")
}

func (rt *Router) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := rt.writeJSON(w, status, env, nil)
	if err != nil {
		rt.logError(err, r, http.StatusInternalServerError)
		w.WriteHeader(500)
	}
}

func (rt *Router) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	rt.logError(err, r, http.StatusInternalServerError)

	m := "internal server error"
	rt.errorResponse(w, r, http.StatusInternalServerError, m)
}

func (rt *Router) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	m := "not found"
	rt.errorResponse(w, r, http.StatusNotFound, m)
	rt.logWarn(r, http.StatusNotFound)
}
