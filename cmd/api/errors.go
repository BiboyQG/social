package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "error", err.Error(), "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
	errorJSON(w, http.StatusInternalServerError, "server encountered an error")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "error", err.Error(), "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
	errorJSON(w, http.StatusNotFound, "resource not found")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "error", err.Error(), "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
	errorJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflict(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "error", err.Error(), "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
	errorJSON(w, http.StatusConflict, err.Error())
}
