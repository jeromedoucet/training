package dao

import (
	"context"
	"database/sql"
	"log"

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
func (u *UserDAO) Insert(ctx context.Context, user *model.User, password string) (*model.User, error) {

	user.Id = uuid.New()

	_, err := u.insertUser.ExecContext(ctx, user.Id.String(), user.Login, user.FirstName, user.LastName, user.Email, password)
	if err != nil {
		return nil, err
	}

	return user, nil

}
