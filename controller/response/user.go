package response

import "github.com/jeromedoucet/training/model"

// User is the type used as a response to
// a sign-in or a login request
type User struct {
	Id        string `id:"json"`
	Login     string `login:"json"`
	Email     string `email:"json"`
	FirstName string `firstName:"json"`
	LastName  string `lastName:"json"`
}

// FromUserModel create a user response from a model
func FromUserModel(user *model.User) *User {
	return &User{Id: user.Id.String(), Login: user.Login, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}
}
