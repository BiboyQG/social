package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.RemoteAddr, r.URL.Path, err.Error())
	errorJSON(w, http.StatusInternalServerError, "server encountered an error")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.RemoteAddr, r.URL.Path, err.Error())
	errorJSON(w, http.StatusNotFound, "resource not found")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.RemoteAddr, r.URL.Path, err.Error())
	errorJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflict(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", r.RemoteAddr, r.URL.Path, err.Error())
	errorJSON(w, http.StatusConflict, err.Error())
}
