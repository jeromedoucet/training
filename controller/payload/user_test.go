package payload_test

import (
	"testing"

	"github.com/jeromedoucet/training/controller/payload"
)

func TestCheckUserCreationPayloadValid(t *testing.T) {
	user := &payload.User{}

	if user.CreationPayloadValid() {
		t.Fatal("Expect user not to be valid")
	}

	user.Login = "jerdct"

	if user.CreationPayloadValid() {
		t.Fatal("Expect user not to be valid")
	}

	user.FirstName = "Jerome"

	if user.CreationPayloadValid() {
		t.Fatal("Expect user not to be valid")
	}

	user.LastName = "Doucet"

	if user.CreationPayloadValid() {
		t.Fatal("Expect user not to be valid")
	}

	user.Email = "jerdct@gmail.com"

	if user.CreationPayloadValid() {
		t.Fatal("Expect user not to be valid")
	}

	user.Password = "password"

	if !user.CreationPayloadValid() {
		t.Fatal("Expect user to be valid")
	}

}
