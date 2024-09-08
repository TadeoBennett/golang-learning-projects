package postgresql

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"tadeobennett.net/quotation/pkg/models"
	// "errors"
	// "tadeobennett.net/quotation/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) //iterates the hash 12  times
	if err != nil {
		return err
	}

	s := `INSERT INTO users(name, email, password_hash)
	VALUES ($1, $2, $3)`

	// INSERT INTO users(name, email, password_hash, activated) VALUES ('Tadeo', 'tadeo@gmail.com', 'tadeo2002', true);

	_, err = m.DB.Exec(s, name, email, hashedPassword) //just returns the row and error

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

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	s := `
	SELECT id, password_hash
	FROM users
	WHERE email = $1
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
			return 0, err //return whatever err it is
		}
	}

	return id, nil
}

// func (m *UserModel) Get(id int) (*models.User, error){
// 	return 0, nil
// }
