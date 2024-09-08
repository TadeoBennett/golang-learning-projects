package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// returns the mux
// makes the routes function a method of application
func (app *application) routes() http.Handler {

	//create a variable to hold my middleware chain in order
	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeadersMiddleware,
	)

	//loads and saves session data to and from the session cookie
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	//pat is a third party library to create and handle routers/multiplexer
	mux := pat.New()
	// Register a catch-all route using http.NotFoundHandler
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/quote/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createQuote)) //added a require authentication for seeing this form
	mux.Post("/quote/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createQuote)) //post request
	mux.Get("/quote/:id", dynamicMiddleware.ThenFunc(app.showQuote))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	//create a file server to serve out static content
	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)

	// return app.recoverPanicMiddleware(app.logRequestMiddleware(securityHeadersMiddleware(mux)))
}
