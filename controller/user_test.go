package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// sign - in
	//  - passwords requirements ?
	// login in
	//  - Nominal
	//  - 401 bad password
	//  - 401 inexisting user
	//  - missing login or password
	// delete

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

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/users", s.URL), bytes.NewBuffer(body))
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

	header := resp.Header.Get("authorization")

	if len(header) == 0 {
		t.Fatal("Expect authorization header not to be empty")
	}

	chunck := strings.Split(strings.TrimSpace(header), " ")
	if len(chunck) != 2 {
		t.Fatal("Expect authorization header to have the form Bearer <TOKEN>")
	}

	var tokenValid bool
	tokenValid, err = test.TokenValid(chunck[1], conf.JwtSecret)

	if err != nil {
		t.Fatalf("Expect no error when validating token, but get %s", err)
	}

	if !tokenValid {
		t.Fatalf("Expect token %s to be valid", chunck[1])
	}

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

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/users", s.URL), bytes.NewBuffer(body))
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

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/users", s.URL), bytes.NewBuffer(body))
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

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/login", s.URL), bytes.NewBuffer(body))
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

	header := resp.Header.Get("authorization")

	if len(header) == 0 {
		t.Fatal("Expect authorization header not to be empty")
	}

	chunck := strings.Split(strings.TrimSpace(header), " ")
	if len(chunck) != 2 {
		t.Fatal("Expect authorization header to have the form Bearer <TOKEN>")
	}

	var tokenValid bool
	tokenValid, err = test.TokenValid(chunck[1], conf.JwtSecret)

	if err != nil {
		t.Fatalf("Expect no error when validating token, but get %s", err)
	}

	if !tokenValid {
		t.Fatalf("Expect token %s to be valid", chunck[1])
	}
}
