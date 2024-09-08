package main

import (
	"net/url"

	"tadeobennett.net/quotation/pkg/models"
)

type templateData struct {
	CSRFToken       string
	Quotes          []*models.Quote
	Quote           *models.Quote
	Flash           string
	ErrorsFromForm  map[string]string //map[key]//value
	FormData        url.Values        //asking go to get the values from the form
	IsAuthenticated bool
}
