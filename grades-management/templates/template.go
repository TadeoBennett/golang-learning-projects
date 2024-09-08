package templates

import (
	"net/url"

	"grademgmt.com/final/pkg/models"
)

type TemplateData struct {
	CSRFToken                string
	Students                 []*models.Student
	Student                  *models.Student
	Grades                   []*models.Grade
	Subjects                 []*models.Subject
	GradeGroupingByStudentId []*models.GradeGroupingByStudentId
	StudentAndAverageGrade   []*models.StudentAndAverageGrade
	Grade                    models.Grade
	ErrorsFromForm           map[string]string
	FormData                 url.Values
	Flash                    string
	IsAuthenticated          bool
}
