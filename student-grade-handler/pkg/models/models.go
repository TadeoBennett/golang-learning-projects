//Cahlil Tillett
//Tadeo Bennett

// need to give the package a main (same as the directory it is in)
package models

import (
	"errors"
	"time"
)

var ( //so we don't have to write var for every variable declaration
	//creating new errors for this model
	ErrRecordNotFound     = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// a struct to hold a Student
type Student struct {
	Student_ID int
	Firstname  string
	Lastname   string
	Email      string
	Address    string
	Age        int
	Created_At time.Time
}

type Teacher struct {
	
}
