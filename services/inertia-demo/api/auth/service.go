package auth

import (
	"database/sql"
	"errors"

	"github.com/bedrock/services/inertia-demo/api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	db *sql.DB
}

var errInvalidCredentials = errors.New("auth: invalid credentials")

func newService(db *sql.DB) (service, error) {
	if db == nil {
		return service{}, errors.New("auth: database connection must not be nil")
	}

	return service{db: db}, nil
}

func (s service) authenticate(form loginForm) (*database.User, error) {
	user, err := database.FindUserByEmail(s.db, form.Email)

	if err != nil {
		return nil, err
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password)) != nil {
		return nil, errInvalidCredentials
	}

	return user, nil
}
