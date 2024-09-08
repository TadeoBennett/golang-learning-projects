package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *AppConfig) routes() http.Handler {
	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeadersMiddleware,
	)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/teacher/signup", dynamicMiddleware.ThenFunc(app.signupTeacherForm)) //added a require authentication for seeing this form
	mux.Post("/teacher/signup", dynamicMiddleware.ThenFunc(app.signupTeacher))   //added a require authentication for seeing this form
	mux.Get("/teacher/login", dynamicMiddleware.ThenFunc(app.loginTeacherForm))   //Post request
	mux.Post("/teacher/login", dynamicMiddleware.ThenFunc(app.loginTeacher))


	mux.Post("/teacher/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutTeacher))
	mux.Get("/forms", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.displayForms))
	mux.Post("/student/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createStudent))
	mux.Get("/student/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getGradesForStudent))
	mux.Post("/grade/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createGradeForStudent))
	mux.Get("/show", dynamicMiddleware.ThenFunc(app.displayGrades))

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}
