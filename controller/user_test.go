package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/controller"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/model"
	"github.com/jeromedoucet/training/test"
)

func TestUserSuite(t *testing.T) {

	test.CleanDB(db)
	t.Run("nominal sign-in", nominalSignIn)

	test.CleanDB(db)
	t.Run("conflict during sign-in", conflictSignIn)

	test.CleanDB(db)
	t.Run("missing fields when sign-in", missingFieldsSignIn)

	test.CleanDB(db)
	t.Run("nominal login", nominalLogIn)

	test.CleanDB(db)
	t.Run("bad identifier login", badIdentifierLogIn)

	test.CleanDB(db)
	t.Run("bad password login", badPasswordLogIn)

	test.CleanDB(db)
	t.Run("missing field login", MissingFieldsLogIn)

	// sign - in
	//  - passwords requirements ?
	// delete

	// todo refresh token ? csrf token ?

}

func nominalSignIn(t *testing.T) {
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Login     string `json:"login"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		"jerdct",
		"Jérôme",
		"Doucet",
		"jerdct@gmail.com",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/users", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 return code. Got %d", resp.StatusCode)
	}

	test.CheckAuthCookie(resp.Cookies(), conf, t)

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	// check some properties of the response body
	var createdUser *response.User
	json.Unmarshal(payloadResp, &createdUser)

	if len(createdUser.Id) == 0 {
		t.Fatal("expect non empty id")
	}
}

func conflictSignIn(t *testing.T) {
	test.InsertUser(&model.User{Id: uuid.New(), Login: "jerdct"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Login     string `json:"login"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		"jerdct",
		"Jérôme",
		"Doucet",
		"jerdct@gmail.com",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/users", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("Expected 409 return code. Got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	// check some properties of the response body
	var res *response.Error
	json.Unmarshal(payloadResp, &res)

	if res.Message != "Another user already exist with this identifier" {
		t.Fatalf("Expect %s, got %s", "Another user already exist with this identifier", res.Message)
	}
}

func missingFieldsSignIn(t *testing.T) {
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		"Jérôme",
		"Doucet",
		"jerdct@gmail.com",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/users", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 return code. Got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	// check some properties of the response body
	var res *response.Error
	json.Unmarshal(payloadResp, &res)

	if res.Message != "Missing some mandatory fields" {
		t.Fatalf("Expect %s, got %s", "Missing some mandatory fields", res.Message)
	}
}

func nominalLogIn(t *testing.T) {
	test.InsertUser(&model.User{Id: uuid.New(), Login: "jerdct", Password: "titi_123456_tata"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		"jerdct",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/login", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 return code. Got %d", resp.StatusCode)
	}

	test.CheckAuthCookie(resp.Cookies(), conf, t)
}

func badIdentifierLogIn(t *testing.T) {
	test.InsertUser(&model.User{Id: uuid.New(), Login: "titi", Password: "titi_123456_tata"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		"jerdct",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/login", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 return code. Got %d", resp.StatusCode)
	}

	header := resp.Header.Get("authorization")

	if len(header) != 0 {
		t.Fatalf("Expect authorization header to be empty but got %s", header)
	}
}

func badPasswordLogIn(t *testing.T) {
	test.InsertUser(&model.User{Id: uuid.New(), Login: "jerdct", Password: "titi"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		"jerdct",
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/login", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 return code. Got %d", resp.StatusCode)
	}

	header := resp.Header.Get("authorization")

	if len(header) != 0 {
		t.Fatalf("Expect authorization header to be empty but got %s", header)
	}
}

func MissingFieldsLogIn(t *testing.T) {
	test.InsertUser(&model.User{Id: uuid.New(), Login: "jerdct", Password: "titi"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Password string `json:"password"`
	}{
		"titi_123456_tata",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/login", s.URL), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 return code. Got %d", resp.StatusCode)
	}

	if test.GetCookieByName(resp.Cookies(), "auth") != nil {
		t.Fatal("Expect to have no auth cookie, but got one")

	}

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	// check some properties of the response body
	var res *response.Error
	json.Unmarshal(payloadResp, &res)

	if res.Message != "Missing some mandatory fields" {
		t.Fatalf("Expect %s, got %s", "Missing some mandatory fields", res.Message)
	}
}
