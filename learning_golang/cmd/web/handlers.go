package main

import (
	// "fmt"

	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/justinas/nosurf"
	"tadeobennett.net/quotation/pkg/models"
)

// displays all the quotes
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("page not found")
		app.notFound(w) //use our custom log
		return
	}
	//the above is taken care of in the routes.go using pat (it only looks for "/")

	q, err := app.quotes.Read()

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	//an instance of template data -------------------------------
	data := &templateData{
		Quotes:          q,
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
	//Display quotes using a template
	ts, err := template.ParseFiles("../../ui/html/show_page.tmpl")

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	//if there are no errors
	err = ts.Execute(w, data)

	if err != nil {
		log.Panicln(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) createQuoteForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/html/quotes_form_page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	//if there are no errors
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createQuote(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//get the values from the request
	author := r.PostForm.Get("author_name")
	category := r.PostForm.Get("category")
	quote := r.PostForm.Get("quote")

	//check the web form fields for validity. We will use a map to save the errrors
	errors := make(map[string]string)

	//check each field
	if strings.TrimSpace(author) == "" {
		errors["author"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(author) > 25 { //RunCountInString is used to count the characters
		errors["author"] = "This field is too long (max 25 characters)"
	}

	if strings.TrimSpace(category) == "" {
		errors["category"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(category) > 15 { //RunCountInString is used to count the characters
		errors["category"] = "This field is too long (max 15 characters)"
	}

	if strings.TrimSpace(quote) == "" {
		errors["quote"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(quote) > 50 { //RunCountInString is used to count the characters
		errors["quote"] = "This field is too long (max 50 characters)"
	}

	if len(errors) > 0 { //an error exists
		ts, err := template.ParseFiles("../../ui/html/quotes_form_page.tmpl") //load the template file

		if err != nil { //error loading the template
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.IsAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			log.Panicln(err.Error())
			app.serverError(w, err)
			return
		}

		return
	}

	//insert a quote
	id, err := app.quotes.Insert(author, category, quote)

	//check if an error was returned from the insert function
	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	//set some session data after a quote is added
	app.session.Put(r, "flash", "Quote Successfully added")

	http.Redirect(w, r, fmt.Sprintf("/quote/%d", id), http.StatusSeeOther)

}

func (app *application) showQuote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi((r.URL.Query().Get(":id"))) //gets the ID from the URL as an integer instead of a string

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	q, err := app.quotes.Get(id)

	//check the errors returned
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	//look for the entry in the session and once found, is deleted from the esession data
	//contains an empty string or a flash message
	flash := app.session.PopString(r, "flash")

	data := &templateData{
		Quote:           q,
		Flash:           flash,
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	// Display the quote using a template
	ts, err := template.ParseFiles("../../ui/html/quote_page.tmpl")
	if err != nil { //error loading the template
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Panicln(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "This is the signup form")
	ts, err := template.ParseFiles("../../ui/html/signup.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	//if there are no errors
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Added a new user")

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//get the values from the request
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	//check the web form fields for validity. We will use a map to save the errrors
	// errors := make(map[typeofKEY]typeofVALUE)
	errors_user := make(map[string]string)

	//check each field
	if strings.TrimSpace(name) == "" {
		errors_user["name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(name) > 50 { //RunCountInString is used to count the characters
		errors_user["name"] = "This field is too long"
	} else if utf8.RuneCountInString(name) < 5 { //RunCountInString is used to count the characters
		errors_user["name"] = "This field is too short"
	}

	if strings.TrimSpace(email) == "" {
		errors_user["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 60 { //RunCountInString is used to count the characters
		errors_user["email"] = "This field is too long"
	}

	//check if an email is properly formed(valid)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		errors_user["email"] = "Invalid Email"
	}

	if strings.TrimSpace(password) == "" {
		errors_user["password"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(password) > 50 { //RunCountInString is used to count the characters
		errors_user["password"] = "This field is too long"
	} else if utf8.RuneCountInString(password) < 8 { //RunCountInString is used to count the characters
		errors_user["password"] = "This field is too short"
	}

	if len(errors_user) > 0 { //an error exists
		ts, err := template.ParseFiles("../../ui/html/signup.page.tmpl")

		if err != nil { //error loading the template
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors_user,
			FormData:        r.PostForm,
			IsAuthenticated: app.IsAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			log.Panicln(err.Error())
			app.serverError(w, err)
			return
		}
		return
	}

	// insert a user
	err = app.users.Insert(name, email, password)
	//check if an error was returned from the insert function
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errors_user["email"] = "Email already in use"
			//redisplay the signup form heading
			ts, err := template.ParseFiles("../../ui/html/signup.page.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			//if there are no errors
			err = ts.Execute(w, &templateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.IsAuthenticated(r),
				CSRFToken:       nosurf.Token(r),
			})
			if err != nil {
				app.serverError(w, err)
			}
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	//set some session data after a quote is added
	app.session.Put(r, "flash", "User Successfully added")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "This is the signup form")
	ts, err := template.ParseFiles("../../ui/html/login.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	//if there are no errors
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	errors_user := make(map[string]string)
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			errors_user["default"] = "Email or Password is Incorrect"
			ts, err := template.ParseFiles("../../ui/html/login.page.tmpl") //load the template file

			if err != nil {
				log.Println(err.Error())
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &templateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.IsAuthenticated(r),
				CSRFToken:       nosurf.Token(r),
			})
			if err != nil {
				log.Panicln(err.Error())
				app.serverError(w, err)
				return
			}
			return
		}
		return
	}
	app.session.Put(r, "authenticatedUserId", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserId")
	app.session.Put(r, "flash", "You have been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther) //go to home when logged out
}
