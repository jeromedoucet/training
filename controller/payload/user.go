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

// ToModel turn a user request payload into a model structure
func (p *User) ToModel() *model.User {
	return &model.User{Login: p.Login, FirstName: p.FirstName, LastName: p.LastName, Email: p.LastName, Password: p.Password}
}
