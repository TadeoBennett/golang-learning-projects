//Cahlil Tillett
//Tadeo Bennett

package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeadersMiddleware,
	)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/student/add", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createStudentForm))
	mux.Post("/student/add", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createStudent))
	mux.Post("/student/find", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getStudent))
	mux.Get("/student/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getStudent)) //Cahlil
	mux.Get("/show", dynamicMiddleware.ThenFunc(app.displayStudents))

	mux.Get("/teacher/signup", dynamicMiddleware.ThenFunc(app.signupTeacherForm)) //Cahlil
	mux.Post("/teacher/signup", dynamicMiddleware.ThenFunc(app.signupTeacher))    //Tadeo
	mux.Get("/teacher/login", dynamicMiddleware.ThenFunc(app.loginTeacherForm))   //Tadeo
	mux.Post("/teacher/login", dynamicMiddleware.ThenFunc(app.loginTeacher))      //Cahlil
	mux.Post("/teacher/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutTeacher))

	//create a file server to serve out static content
	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	//the initial request will be passed to the log
	return standardMiddleware.Then(mux)
}
