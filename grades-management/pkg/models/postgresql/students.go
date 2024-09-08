package postgresql

import (
	"database/sql"

	"grademgmt.com/final/pkg/models"
)

type StudentModel struct {
	DB *sql.DB
}

func (m *StudentModel) InsertStudent(student_id, firstname, lastname string, age int, gender string) (int, error) {
	var id int

	s := `
    INSERT INTO students (student_id, firstname, lastname, age, gender)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING student_id;
    `

	err := m.DB.QueryRow(s, student_id, firstname, lastname, age, gender).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *StudentModel) ReadAllStudents() ([]*models.Student, error) {
	s := `
	SELECT * FROM students;
	`

	rows, err := m.DB.Query(s) //returns result rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []*models.Student{}

	for rows.Next() {
		st := &models.Student{}
		err = rows.Scan(&st.ID, &st.Student_ID, &st.Firstname, &st.Lastname, &st.Age, &st.Gender, &st.Created_At)
		if err != nil {
			//the slice is empty
			return nil, err
		}
		students = append(students, st)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (m *StudentModel) ReadAllSubjects() ([]*models.Subject, error) {
	s := `
	SELECT subject FROM grades
	GROUP BY subject;
	`

	rows, err := m.DB.Query(s) //returns result rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := []*models.Subject{}

	for rows.Next() {
		sj := &models.Subject{}
		err = rows.Scan(&sj.Subject)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, sj)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subjects, nil
}

func (m *StudentModel) ReadStudentWithId(id int) (*models.Student, error) {
	s := `
    SELECT *
    FROM students 
    WHERE student_id = ?;
    `
	row := m.DB.QueryRow(s, id)

	// Create a new student object to store the result
	student := &models.Student{}

	// Scan the row into the student object
	err := row.Scan(&student.ID, &student.Student_ID, &student.Firstname, &student.Lastname, &student.Age, &student.Gender, &student.Created_At)
	if err != nil {
		// Handle the error, which could be sql.ErrNoRows if no rows were found
		return nil, err
	}

	// Return the student object
	return student, nil
}
