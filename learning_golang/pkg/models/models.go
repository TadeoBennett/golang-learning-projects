// need to give the package a main (same as the directory it is in)
package models

import (
	"errors"
	"time"
)

var( //so we don't have to write var for every variable declaration
	//creating new errors for this model
	ErrRecordNotFound = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//a struct to hold a quote
type Quote struct {
	Quotation_id   int
	Created_at time.Time
	Author_name    string
	Category       string
	Body          string
}

type User struct{
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
	Active bool
}

