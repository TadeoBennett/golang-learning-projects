//  Cahlil Tillett
//  Tadeo Bennett 
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

func securityHeadersMiddleware(next http.Handler) http.Handler {
	// note: middleware has to return a handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//added 2 security headers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		//call the original handler after adding the two headers //continue the chain
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	// note: middleware has to return a handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Preprocessing----------------------------------------------------
		start := time.Now()
		//Continue the chain--------------------------------------------
		next.ServeHTTP(w, r)
		//Postprocessing---------------------------------------------
		app.infoLog.Printf("%s %s %s %s %s",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func (app *application) recoverPanicMiddleware(next http.Handler) http.Handler {
	// note: middleware has to return a handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//function is only called if it encounters a panic
		defer func() {
			if err := recover(); err != nil { //if err is not nil
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}() //this means that it will execute
		next.ServeHTTP(w, r)
	})
}


func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.session.Put(r, "flash", "You need to login to perform that action.")
			http.Redirect(w, r, "/teacher/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
