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

	"github.com/justinas/nosurf"
	"grademgmt.com/final/pkg/models"
	"grademgmt.com/final/templates"
)

func (app *AppConfig) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/index.html")

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash")

	err = ts.Execute(w, &templates.TemplateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		log.Panicln(err.Error())
		app.serverError(w, err)
	}
	log.Println("Homepage displayed")
}
func (app *AppConfig) displayForms(w http.ResponseWriter, r *http.Request) {
	s, err := app.students.ReadAllStudents() ///gets the list of students from the database

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}
	sj, err := app.students.ReadAllSubjects() ///gets the list of students from the database

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash")

	//an instance of template data -------------------------------
	data := &templates.TemplateData{
		Flash:           flash,
		Students:        s,
		Subjects:        sj,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	//Display Students using a template
	ts, err := template.ParseFiles("../../ui/tmpl/form.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	//if there are no errors
	err = ts.Execute(w, data)

	if err != nil {
		panic(err)
	}

	log.Println("Form page displayed. Loaded all students and subjects to necessary inputs")
}

func (app *AppConfig) signupTeacherForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/signup.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, &templates.TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *AppConfig) signupTeacher(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	firstname := strings.TrimSpace(r.PostForm.Get("firstname"))
	lastname := strings.TrimSpace(r.PostForm.Get("lastname"))
	email := strings.TrimSpace(r.PostForm.Get("email"))
	password := strings.TrimSpace(r.PostForm.Get("password"))

	errors_teacher := make(map[string]string)

	//check each field
	if firstname == "" {
		errors_teacher["fname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(firstname) > 50 {
		errors_teacher["fname"] = "This field is too long"
	} else if utf8.RuneCountInString(firstname) < 3 {
		errors_teacher["fname"] = "This field is too short"
	}

	if lastname == "" {
		errors_teacher["lname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(firstname) > 50 {
		errors_teacher["lname"] = "This field is too long"
	} else if utf8.RuneCountInString(firstname) < 3 {
		errors_teacher["lname"] = "This field is too short"
	}

	if strings.TrimSpace(email) == "" {
		errors_teacher["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 60 { //RunCountInString is used to count the characters
		errors_teacher["email"] = "This field is too long"
	}

	//check if an email is properly formed(valid)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		errors_teacher["email"] = "Invalid Email"
	}

	if password == "" {
		errors_teacher["password"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(password) > 50 {
		errors_teacher["password"] = "This field is too long"
	} else if utf8.RuneCountInString(password) < 8 {
		errors_teacher["password"] = "This field is too short"
	}

	if len(errors_teacher) > 0 { //an error exists
		ts, err := template.ParseFiles("../../ui/tmpl/signup.tmpl")

		if err != nil { //error loading the template
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, &templates.TemplateData{
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

	// insert a user
	err = app.teachers.Insert(firstname, lastname, email, password)
	//check if an error was returned from the insert function
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errors_teacher["email"] = "Email already in use"
			//redisplay the signup form heading
			ts, err := template.ParseFiles("../../ui/tmpl/signup.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			//if there are no errors
			err = ts.Execute(w, &templates.TemplateData{
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
	app.session.Put(r, "flash", "Teacher Successfully added")

	http.Redirect(w, r, "/teacher/login", http.StatusSeeOther)
}

func (app *AppConfig) loginTeacherForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../../ui/tmpl/login.tmpl")

	if err != nil {
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash")

	err = ts.Execute(w, &templates.TemplateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *AppConfig) createStudent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fname := strings.TrimSpace(r.PostForm.Get("firstname"))
	lname := strings.TrimSpace(r.PostForm.Get("lastname"))
	studentid := strings.TrimSpace(r.PostForm.Get("studentid"))
	age := r.PostForm.Get("age")
	gender := r.PostForm.Get("gender")

	errors_student := make(map[string]string)

	if studentid == "" {
		errors_student["studentid"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(studentid) != 10 {
		errors_student["studentid"] = "Invalid ID. Student IDs are 10 chars long"
	} else if !containsOnlyValidChars(fname) {
		errors_student["studentid"] = "Please provide valid characters"
	}

	if fname == "" {
		errors_student["fname"] = "Please provide a firstname"
	} else if utf8.RuneCountInString(fname) > 20 {
		errors_student["fname"] = "This field is too long"
	} else if !containsOnlyValidChars(fname) {
		errors_student["fname"] = "Please provide valid characters"
	}
	if lname == "" {
		errors_student["lname"] = "Please provide a lastname"
	} else if utf8.RuneCountInString(lname) > 20 {
		errors_student["lname"] = "This field is too long"
	} else if !containsOnlyValidChars(lname) {
		errors_student["lname"] = "Please provide valid characters"
	}

	inputAgeInt, err := strconv.Atoi(strings.TrimSpace(age))
	if err != nil {
		errors_student["age"] = "Internal Error"
	}
	if age == "" {
		errors_student["age"] = "This field cannot be left blank"
	} else if inputAgeInt > 99 {
		errors_student["age"] = "Please provide a valid number"
	} else if utf8.RuneCountInString(age) >= 3 {
		errors_student["age"] = "This field is too long"
	}

	if gender == "" {
		errors_student["gender"] = "Please provide a gender"
	} else if utf8.RuneCountInString(gender) > 1 {
		errors_student["gender"] = "Please select a gender"
	} else if gender != "M" && gender != "F" {
		errors_student["gender"] = "Invalid gender provided2"
	}

	if len(errors_student) > 0 {
		log.Println("Error found in form inputs")
		students, err := app.students.ReadAllStudents()
		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}
		ts, err := template.ParseFiles("../../ui/tmpl/form.tmpl")
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, &templates.TemplateData{
			Students:        students,
			ErrorsFromForm:  errors_student,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			panic(err)
		}
		return
	}

	id, err := app.students.InsertStudent(studentid, fname, lname, inputAgeInt, gender)

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Student Successfully added")

	http.Redirect(w, r, fmt.Sprintf("/forms?newstudentid=%d", id), http.StatusSeeOther)
}

func (app *AppConfig) getGradesForStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi((r.URL.Query().Get(":id"))) //gets the ID from the URL as an integer instead of a string

	if err != nil || id < 1 {
		app.errRecordNotFound(w, err)
		return
	}
	s, err := app.grades.ReadStudentAndAverageGrade(id)

	if err != nil {
		log.Println("no grades for student with id", id)
		app.serverError(w, err)
		return
	}

	flash := app.session.PopString(r, "flash")
	// flash = "Hello there"

	data := &templates.TemplateData{
		Flash:                  flash,
		StudentAndAverageGrade: s,
		IsAuthenticated:        app.isAuthenticated(r),
		CSRFToken:              nosurf.Token(r),
	}

	// Display the quote using a template
	ts, err := template.ParseFiles("../../ui/tmpl/student.tmpl")
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

func (app *AppConfig) createGradeForStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid method. Expected Post")
		http.Redirect(w, r, "/forms?invalidmethod", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	studentid := r.PostForm.Get("studentidselect")
	subject := strings.TrimSpace(r.PostForm.Get("createdSubject"))
	grade := r.PostForm.Get("grade")
	errors := make(map[string]string)

	inputStudentIDInt, err := strconv.Atoi(strings.TrimSpace(studentid))
	if err != nil {
		errors["studentidselect"] = "Internal Error"
	}
	if studentid == "" {
		errors["studentidselect"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(studentid) < 0 {
		errors["studentidselect"] = "Please select a valid student"
	}

	if subject == "" {
		errors["subject"] = "Please provide a subject"
	} else if utf8.RuneCountInString(subject) > 30 {
		errors["subject"] = "This field is too long"
	} else if !containsOnlyValidChars(subject) {
		errors["subject"] = "Please provide valid characters"
	}

	inputGradeInt, err := strconv.Atoi(strings.TrimSpace(grade))
	if err != nil {
		errors["grade"] = "Internal Error"
	}
	if grade == "" {
		errors["grade"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(grade) >= 4 {
		errors["grade"] = "This field is too long"
	} else if inputGradeInt > 100 {
		errors["grade"] = "Please provide a valid number(<=100)"
	}

	if len(errors) > 0 {
		log.Println("Error found in form inputs")
		log.Println(errors)

		students, err := app.students.ReadAllStudents() ///gets the list of students from the database
		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}
		ts, err := template.ParseFiles("../../ui/tmpl/form.tmpl")
		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, &templates.TemplateData{
			Students:        students,
			ErrorsFromForm:  errors,
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

	id, err := app.grades.InsertGrade(inputStudentIDInt, subject, inputGradeInt)

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Grade Successfully added")

	url := fmt.Sprintf("/forms?new_grade_id=%d", id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *AppConfig) displayGrades(w http.ResponseWriter, r *http.Request) {
	g, err := app.grades.ReadGradesGroupByStudent()

	if err != nil {
		log.Println(err.Error())
		app.serverError(w, err)
		return
	}

	//an instance of template data -------------------------------
	data := &templates.TemplateData{
		GradeGroupingByStudentId: g,
		IsAuthenticated:          app.isAuthenticated(r),
		CSRFToken:                nosurf.Token(r),
	}

	//Display Students using a template
	ts, err := template.ParseFiles("../../ui/tmpl/list.tmpl")

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

	log.Println("Showing Grouping of all students' grades")
}

func (app *AppConfig) loginTeacher(w http.ResponseWriter, r *http.Request) {

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
		if errors.Is(err, models.ErrInvalidCredentials) {
			errors_teacher["default"] = "Email or Password is Incorrect"
			ts, err := template.ParseFiles("../../ui/tmpl/login.tmpl") //load the template file

			if err != nil {
				log.Println(err.Error())
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &templates.TemplateData{
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

func (app *AppConfig) logoutTeacher(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedTeacherId")
	app.session.Put(r, "flash", "You have been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther) //go to home when logged out
}

