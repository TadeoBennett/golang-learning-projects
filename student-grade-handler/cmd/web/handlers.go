//Cahlil Tillett
//Tadeo Bennett

package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"advancedweb.com/test2/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
	}
	ts, err := template.ParseFiles("../../ui/tmpl/index.tmpl")

	if err != nil {
		panic(err.Error())
	}

	flash := app.session.PopString(r, "flash")

	err = ts.Execute(w, &TemplateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		panic(err.Error())
	}
}

func (app *application) createStudentForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/studentform.page.tmpl")

	if err != nil {
		// panic(err.Error())
		app.serverError(w, err)
		// return
	}

	err = ts.Execute(w, &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		panic(err.Error())
	}

	log.Println("showing the student form page")
}

func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid method access to /student/add; redirecting...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	fname := r.PostForm.Get("first_name")
	lname := r.PostForm.Get("last_name")
	email := r.PostForm.Get("email")
	address := r.PostForm.Get("address")
	age := r.PostForm.Get("age")

	errors_student := make(map[string]string)

	if strings.TrimSpace(fname) == "" {
		errors_student["fname"] = "This field cannot be left blank"
	}else if(!containsOnlyValidChars("fname")){
		errors_student["fname"] = "This field contains invalid characters"
	} else if utf8.RuneCountInString(fname) > 25 {
		errors_student["fname"] = "This field is too long"
	}

	if strings.TrimSpace(lname) == "" {
		errors_student["lname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(lname) > 25 {
		errors_student["lname"] = "This field is too long"
	}

	if strings.TrimSpace(email) == "" {
		errors_student["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 50 {
		errors_student["email"] = "This field is too long"
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		errors_student["email"] = "Invalid Email"
	}

	if strings.TrimSpace(address) == "" {
		errors_student["address"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(address) > 50 {
		errors_student["address"] = "This field is too long"
	}

	inputAgeInt, err := strconv.Atoi(strings.TrimSpace(age))

	if strings.TrimSpace(age) == "" {
		errors_student["age"] = "This field cannot be left blank"
	} else if inputAgeInt > 100 || err != nil {
		errors_student["age"] = "This field cannot contain an invalid number"
	} else if utf8.RuneCountInString(age) > 3 {
		errors_student["age"] = "This field is too long"
	}

	if len(errors_student) > 0 {
		log.Println("Errors in form")
		ts, err := template.ParseFiles("../../ui/tmpl/studentform.page.tmpl")

		if err != nil {
			log.Panicln(err.Error())
			app.serverError(w, err)
		}

		err = ts.Execute(w, &TemplateData{
			ErrorsFromForm:  errors_student,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			log.Panicln(err.Error())
			app.serverError(w, err)
		}

		return
	}

	id, err := app.student_details.Insert(fname, lname, email, address, age)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errors_student["email"] = "Email already in use"
			//redisplay the signup form heading
			ts, err := template.ParseFiles("../../ui/tmpl/signup.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			//if there are no errors
			err = ts.Execute(w, &TemplateData{
				ErrorsFromForm:  errors_student,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
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

	log.Println("Added user with id: ", id)
	app.session.Put(r, "flash", "You added a new student")
	http.Redirect(w, r, fmt.Sprintf("/student/%v", id), http.StatusSeeOther)
}

func (app *application) displayStudents(w http.ResponseWriter, r *http.Request) {
	s, err := app.student_details.Read() ///gets the list of students from the database

	if err != nil {
		panic(err.Error())
	}

	flash := app.session.PopString(r, "flash")

	//an instance of template data -------------------------------
	data := &TemplateData{
		Students:        s,
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	//Display Students using a template
	ts, err := template.ParseFiles("../../ui/tmpl/show.page.tmpl")

	if err != nil {
		panic(err.Error())
	}

	//if there are no errors
	err = ts.Execute(w, data)

	if err != nil {
		panic(err.Error())
	}

	log.Println("Loaded all students' data to the page.")
}

func (app *application) signupTeacher(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	firstname := strings.TrimSpace(r.PostForm.Get("firstname"))
	lastname := strings.TrimSpace(r.PostForm.Get("lastname"))
	address := strings.TrimSpace(r.PostForm.Get("address"))
	log.Println(address)
	age := r.PostForm.Get("age")
	email := strings.TrimSpace(r.PostForm.Get("email"))
	password := strings.TrimSpace(r.PostForm.Get("password"))

	errors_teacher := make(map[string]string)

	//validating the firstname
	if firstname == "" {
		errors_teacher["fname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(firstname) > 30 {
		errors_teacher["fname"] = "This field is too long"
	} else if utf8.RuneCountInString(firstname) < 3 {
		errors_teacher["fname"] = "This field is too short"
	}

	//validating the lastname
	if lastname == "" {
		errors_teacher["lname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(firstname) > 30 {
		errors_teacher["lname"] = "This field is too long"
	} else if utf8.RuneCountInString(firstname) < 3 {
		errors_teacher["lname"] = "This field is too short"
	}

	//validating the address
	if address == "" {
		errors_teacher["address"] = "This field cannot be left blank"
	} else if address != "Corozal" && address != "Orange Walk" && address != "Belize" && address != "Cayo" && address != "Stann Creek" && address != "Punta Gorda" {
		errors_teacher["address"] = `Choose an option from the provided fields`
	} else if utf8.RuneCountInString(address) > 15 {
		errors_teacher["address"] = "This field is too long. Choose an option from the provided fields"
	}

	//validating the age input
	ageInt, err := strconv.Atoi(strings.TrimSpace(age))
	if err != nil {
		errors_teacher["age"] = "Internal Error"
	}
	if age == "" {
		errors_teacher["age"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(age) >= 4 {
		errors_teacher["age"] = "This field is too long"
	} else if ageInt > 100 {
		errors_teacher["age"] = "Please provide a valid age"
	}

	//validating the email
	if strings.TrimSpace(email) == "" {
		errors_teacher["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 70 {
		errors_teacher["email"] = "This field is too long"
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		errors_teacher["email"] = "Invalid Email"
	}

	//validating the password
	if password == "" {
		errors_teacher["password"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(password) > 60 {
		errors_teacher["password"] = "This field is too long"
	} else if utf8.RuneCountInString(password) < 5 {
		errors_teacher["password"] = "This field is too short"
	}

	if len(errors_teacher) > 0 { //an error exists
		ts, err := template.ParseFiles("../../ui/tmpl/signup.page.tmpl")

		if err != nil { //error loading the template
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, &TemplateData{
			ErrorsFromForm:  errors_teacher,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			panic(err.Error())
		}
		return
	}

	// insert a user
	err = app.teachers.Insert(firstname, lastname, address, ageInt, email, password)
	//check if an error was returned from the insert function
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errors_teacher["email"] = "Email already exists in the system"
			//redisplay the signup form heading
			ts, err := template.ParseFiles("../../ui/tmpl/signup.page.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			//if there are no errors
			err = ts.Execute(w, &TemplateData{
				ErrorsFromForm:  errors_teacher,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
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
	app.session.Put(r, "flash", "New teacher has been added")

	http.Redirect(w, r, "/teacher/login", http.StatusSeeOther)
}

func (app *application) signupTeacherForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/signup.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) loginTeacherForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/login.page.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash")

	err = ts.Execute(w, &TemplateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) loginTeacher(w http.ResponseWriter, r *http.Request) {
	log.Print("login teacher reached")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	errors_teacher := make(map[string]string)
	id, err := app.teachers.Authenticate(email, password)
	if err != nil {
		log.Print("didnt authentocate")
		if errors.Is(err, models.ErrInvalidCredentials) {
			errors_teacher["default"] = "Email or Password is Incorrect"
			ts, err := template.ParseFiles("../../ui/tmpl/login.page.tmpl") //load the template file
			if err != nil {
				log.Println(err.Error())
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &TemplateData{
				ErrorsFromForm:  errors_teacher,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
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
	app.session.Put(r, "authenticatedTeacherId", id)
	app.session.Put(r, "flash", "You have logged in.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutTeacher(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedTeacherId")
	app.session.Put(r, "flash", "You have logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther) //go to home when logged out
}

func (app *application) getStudent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Retrieve the value of 'student_id' from the form
	studentID := r.FormValue("student_id")
	var id int

	// Check if studentID is not empty
	if studentID != "" {
		// Convert 'studentID' to an integer
		var err error
		id, err = strconv.Atoi(studentID)
		if err != nil {
			// Handle invalid student ID
			app.clientError(w, http.StatusBadRequest)
			return
		}
	} else {
		// Retrieve the 'id' value from the URL
		var err error
		id, err = strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil || id < 1 {
			app.errRecordNotFound(w, err)
			return
		}
	}

	s, err := app.student_details.ReadByID(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash") //displaying that a new student was added

	data := &TemplateData{
		Flash:           flash,
		Student:         s,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	// Display the quote using a template
	ts, err := template.ParseFiles("../../ui/tmpl/student.page.tmpl")
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
