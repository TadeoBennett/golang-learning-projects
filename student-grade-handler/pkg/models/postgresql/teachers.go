//Cahlil Tillett
//Tadeo Bennett

package postgresql

import (
	"database/sql"
	"errors"

	"advancedweb.com/test2/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type TeacherModel struct {
	DB *sql.DB
}

func (m *TeacherModel) Insert(firstname, lastname, address string, ageInt int, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) //iterates the hash 12  times
	if err != nil {
		return err
	}

	s := `INSERT INTO teachers(First_name, Last_Name, Address, Age, Email, Password)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = m.DB.Exec(s, firstname, lastname, address, ageInt, email, hashedPassword)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicated key value violates unique constraint "users_email_key"`:
			return models.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m *TeacherModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	s := `
	SELECT Teacher_ID, Password
	FROM teachers
	WHERE Email = $1
	AND activated = TRUE
	`

	err := m.DB.QueryRow(s, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	//there was no err; check the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
