package model

import (
	"github.com/google/uuid"
)

// User represent an ... hmm user of training.
// note that for security purposed, the password is not
// strored here.
type User struct {
	Id        uuid.UUID
	Login     string
	FirstName string
	LastName  string
	Email     string
}
