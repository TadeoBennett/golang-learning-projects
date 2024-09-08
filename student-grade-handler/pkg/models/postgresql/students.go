//Cahlil Tillett
//Tadeo Bennett

package postgresql

import (
	"database/sql"

	"advancedweb.com/test2/pkg/models"
)

type StudentModel struct {
	DB *sql.DB
}

func (m *StudentModel) Insert(fname, lname, email, address, age string) (int, error) {
	var id int

	s := `
	INSERT INTO students(First_Name, Last_Name, Address, Email, Age)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING Student_ID		
	`

	err := m.DB.QueryRow(s, fname, lname, address, email, age).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *StudentModel) Read() ([]*models.Student, error) {

	s := `
	SELECT Student_ID, First_Name, Last_Name, Address, Email, Age, Created_At
	FROM students
	`

	rows, err := m.DB.Query(s) //returns the rows of results
	if err != nil {
		//because the slice was empty
		return nil, err
	}
	//clean up before leave Read()
	defer rows.Close()

	//save as array of student structs
	students := []*models.Student{}

	//Iterate over rows (a result set)
	for rows.Next() {

		//has to be initialized to empty
		st := &models.Student{}

		err = rows.Scan(&st.Student_ID, &st.Firstname, &st.Lastname, &st.Address, &st.Email, &st.Age, &st.Created_At)
		if err != nil {
			//the slice is empty
			return nil, err
		}

		//Append to students
		students = append(students, st)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (m *StudentModel) ReadByID(id int) (*models.Student, error) {
    s := `
    SELECT Student_ID, First_Name, Last_Name, Address, Email, Age, Created_At
    FROM students
    WHERE Student_ID = $1
    `

    row := m.DB.QueryRow(s, id)

    st := &models.Student{}
    err := row.Scan(&st.Student_ID, &st.Firstname, &st.Lastname, &st.Address, &st.Email, &st.Age, &st.Created_At)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, models.ErrRecordNotFound
        }
        return nil, err
    }

    return st, nil
}