package postgresql

import (
	"database/sql"
	"log"

	"grademgmt.com/final/pkg/models"
)

type GradeModel struct {
	DB *sql.DB
}

func (m *GradeModel) InsertGrade(id int, subject string, grade int) (int, error) {
	s := `
	INSERT INTO grades (student_id, subject, grade)
	VALUES ($1, $2, $3) 
	RETURNING grade_id;
	
	`

	err := m.DB.QueryRow(s, id, subject, grade).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *GradeModel) ReadAllGrades() ([]*models.Grade, error) {

	s := `
	SELECT * FROM grades;
	`

	rows, err := m.DB.Query(s) //returns the rows of results
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []*models.Grade{}
	for rows.Next() { //parses and save the records into a slice
		g := &models.Grade{}

		err = rows.Scan(&g.ID, &g.Student_ID, &g.Subject, &g.Grade, &g.Created_At)
		if err != nil {
			return nil, err
		}

		//Append to grades
		grades = append(grades, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

func (m *GradeModel) ReadGradesGroupByStudent() ([]*models.GradeGroupingByStudentId, error) {
	//groups the grades gotten from a specific student
	s := `
	SELECT s.id AS student_id,
       CONCAT(s.firstname, ' ', s.lastname) AS full_name,
       g.grade_id,
       g.subject,
       g.grade
	FROM students s
	JOIN grades g ON s.id = g.student_id;
	`

	rows, err := m.DB.Query(s) //returns the rows of results
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []*models.GradeGroupingByStudentId{}
	for rows.Next() { //parses and save the records into a slice
		g := &models.GradeGroupingByStudentId{}

		err = rows.Scan(&g.StudentID, &g.FullName, &g.GradeID, &g.Subject, &g.Grade)
		if err != nil {
			return nil, err
		}

		//Append to grades
		grades = append(grades, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

func (m *GradeModel) ReadStudentAndAverageGrade(id int) ([]*models.StudentAndAverageGrade, error) {
	s := `
	SELECT
    s.id AS student_id,
    s.student_id,
    s.firstname,
    s.lastname,
    s.age,
    s.gender,
    g.subject,
    ROUND(AVG(g.grade), 2) AS average_grade
	FROM
		students s
	JOIN
		grades g ON s.id = g.student_id
	WHERE s.id = $1
	GROUP BY
		s.id, s.student_id, s.firstname, s.lastname, s.age, s.gender, g.subject
	ORDER BY
		s.id;
	`

	rows, err := m.DB.Query(s, id)
	if err != nil {
		log.Println("row was empty. could not get grades for student with provided id")
		return nil, err //row was empty
	}
	defer rows.Close()

	studentAverageGrades := []*models.StudentAndAverageGrade{}
	for rows.Next() {
		g := &models.StudentAndAverageGrade{}

		err := rows.Scan(&g.ID, &g.StudentID, &g.Firstname, &g.Lastname, &g.Age, &g.Gender, &g.Subject, &g.AverageGrade)
		if err != nil {
			log.Println("Slice is empty")
			return nil, err //slice was empty
		}
		studentAverageGrades = append(studentAverageGrades, g)
	}

	// Check if no rows were returned
	if len(studentAverageGrades) == 0 {
		return nil, models.ErrNoRecordReturned
	}

	//error iterating the rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return studentAverageGrades, nil
}
