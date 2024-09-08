package main

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
)

func containsOnlyValidChars(input string) bool {
	validChars := regexp.MustCompile(`^[a-zA-Z ']+$`)
	return validChars.MatchString(input)
}

func (app *AppConfig) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *AppConfig) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *AppConfig) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
func (app *AppConfig) errRecordNotFound(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	app.clientError(w, http.StatusInternalServerError)
}

func (app *AppConfig) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedTeacherId")
}
