package models

import (
	"errors"
	"time"
)

var ( //so we don't repeat "var"
	ErrRecordNotFound     = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrNoRecordReturned    = errors.New("models: no record returned")
)

type Student struct {
	ID         int
	Student_ID int
	Firstname  string
	Lastname   string
	Age        int
	Gender     string
	Created_At time.Time
}

type Subject struct {
	Subject string
}

type Grade struct {
	ID         int
	Student_ID int
	Subject    string
	Grade      int
	Created_At time.Time
}

type GradeGroupingByStudentId struct {
	StudentID int
	FullName  string
	GradeID   int
	Subject   string
	Grade     int
}

type StudentAndAverageGrade struct{
	ID int
	StudentID int
	Firstname string
	Lastname string
	Age int
	Gender string
	Subject string
	AverageGrade float32
}

