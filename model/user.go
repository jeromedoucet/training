package model

import (
	"github.com/google/uuid"
)

// User represent an ... hmm user of training.
type User struct {
	Id        uuid.UUID
	Login     string
	FirstName string
	LastName  string
	Email     string
	Password  string
}
