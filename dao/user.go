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
	db            *sql.DB
	insertUser    *sql.Stmt
	getUser       *sql.Stmt
	checkPassword *sql.Stmt
}

func createUserDAO(db *sql.DB) *UserDAO {
	var err error
	u := &UserDAO{db: db}

	if u.insertUser, err = u.db.Prepare(`
		INSERT INTO "user" ("id", "login", "first_name", "last_name", "email", "password") VALUES ($1, $2, $3, $4, $5, crypt($6, gen_salt('bf')))
	`); err != nil {
		log.Fatalf("An error is returned during insertUser statement initialization %s", err.Error())
	}

	if u.getUser, err = u.db.Prepare(`
		SELECT "login", "first_name", "last_name", "email" FROM "user" WHERE "id" = $1
	`); err != nil {
		log.Fatalf("An error is returned during checkPassword statement initialization %s", err.Error())
	}

	if u.checkPassword, err = u.db.Prepare(`
		SELECT "id" FROM "user" WHERE "login" = $1 AND "password" = crypt($2, "password")
	`); err != nil {
		log.Fatalf("An error is returned during checkPassword statement initialization %s", err.Error())
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

// CheckPassword is used for login / password autentication
func (u *UserDAO) CheckPassword(ctx context.Context, user *model.User) (*model.User, *DbError) {
	rows, err := u.checkPassword.QueryContext(ctx, user.Login, user.Password)
	if err != nil {
		return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
	}
	defer rows.Close()
	if rows.Next() {
		var id string
		rows.Scan(&id)
		return u.GetUser(ctx, id)
	} else {
		return nil, &DbError{Message: "authentication failed", Type: NOT_FOUND}
	}
}

// CheckPassword is used for login / password autentication
func (u *UserDAO) GetUser(ctx context.Context, id string) (*model.User, *DbError) {
	var user *model.User
	var userId uuid.UUID
	rows, err := u.getUser.QueryContext(ctx, id)
	if err != nil {
		return user, &DbError{Message: err.Error(), Type: UNKNOWN}
	}

	userId, err = uuid.Parse(id)
	if err != nil {
		return user, &DbError{Message: err.Error(), Type: UNKNOWN}
	}

	defer rows.Close()
	if rows.Next() {
		user = new(model.User)
		user.Id = userId
		rows.Scan(
			&user.Login,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		)
	} else {
		return user, &DbError{Message: "user not found", Type: NOT_FOUND}
	}
	return user, nil
}
