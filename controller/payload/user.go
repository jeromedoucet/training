package payload

import "github.com/jeromedoucet/training/model"

// User represents a payload for request
// on user(s)
type User struct {
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p *User) CreationPayloadValid() bool {
	if len(p.Login) == 0 {
		return false
	}

	if len(p.FirstName) == 0 {
		return false
	}

	if len(p.LastName) == 0 {
		return false
	}

	if len(p.Email) == 0 {
		return false
	}

	if len(p.Password) == 0 {
		return false
	}

	return true
}

func (p *User) AuthenticationPayloadValid() bool {

	if len(p.Login) == 0 {
		return false
	}

	if len(p.Password) == 0 {
		return false
	}

	return true
}

// ToModel turn a user request payload into a model structure
func (p *User) ToModel() *model.User {
	return &model.User{Login: p.Login, FirstName: p.FirstName, LastName: p.LastName, Email: p.Email, Password: p.Password}
}
