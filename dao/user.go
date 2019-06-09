package dao

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

// UserDAO allow accessing user table
type UserDAO struct {
	db         *sql.DB
	insertUser *sql.Stmt
}

func createUserDAO(db *sql.DB) *UserDAO {
	var err error
	u := &UserDAO{db: db}

	if u.insertUser, err = u.db.Prepare(`
		INSERT INTO "user" ("id", "login", "first_name", "last_name", "email", "password") VALUES ($1, $2, $3, $4, $5, crypt($6, gen_salt('bf')))
	`); err != nil {
		log.Fatalf("An error is returned during insertUser statement initialization %s", err.Error())
	}

	return u
}

// Insert a new user
func (u *UserDAO) Insert(ctx context.Context, user *model.User) (*model.User, *DbError) {

	user.Id = uuid.New()

	_, err := u.insertUser.ExecContext(ctx, user.Id.String(), user.Login, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "user_login_key") {
			return nil, &DbError{Message: "Another user already exist with this identifier", Type: CONFLICT}
		} else {
			return nil, &DbError{Message: errMsg, Type: UNKNOWN}
		}
	}

	return user, nil

}
