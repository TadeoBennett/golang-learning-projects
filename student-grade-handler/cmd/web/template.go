//Cahlil Tillett
//Tadeo Bennett

package main

import (
	"net/url"

	"advancedweb.com/test2/pkg/models"
)

type TemplateData struct {
	CSRFToken       string
	Flash           string
	Students        []*models.Student
	Student         *models.Student
	Teachers        []*models.Teacher
	Teacher         *models.Teacher
	ErrorsFromForm  map[string]string //map[key]//value
	FormData        url.Values
	IsAuthenticated bool
}
