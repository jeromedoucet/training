package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	test.CleanDB(db)
	t.Run("logout", logout)

	test.CleanDB(db)
	t.Run("nominal auth check", nominalAuthCheck)

	test.CleanDB(db)
	t.Run("No token auth check", notTokenAuthCheck)

	test.CleanDB(db)
	t.Run("outdated token auth check", outDatedTokenAuthCheck)

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

	token := test.CheckAuthCookie(resp.Cookies(), conf, t)

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	// check some properties of the response body
	var createdUser *response.User
	json.Unmarshal(payloadResp, &createdUser)

	if len(createdUser.Id) == 0 {
		t.Fatal("expect non empty id")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["sub"] != createdUser.Id {
		t.Fatalf("Expected %s. Got %s", createdUser.Id, claims["sub"])
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
	userId := uuid.New()
	test.InsertUser(&model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}, db)
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

	token := test.CheckAuthCookie(resp.Cookies(), conf, t)

	claims := token.Claims.(jwt.MapClaims)

	if claims["sub"] != userId.String() {
		t.Fatalf("Expected %s. Got %s", userId.String(), claims["sub"])
	}

	defer resp.Body.Close()
	payloadResp, _ := ioutil.ReadAll(resp.Body)

	var connectedUser *response.User
	json.Unmarshal(payloadResp, &connectedUser)

	if connectedUser.Id != userId.String() {
		t.Fatalf("expect %s, got %s", userId.String(), connectedUser.Id)
	}

	if connectedUser.Login != "jerdct" {
		t.Fatalf("expect %s, got %s", "jerdct", connectedUser.Login)
	}
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

func logout(t *testing.T) {
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	c := &http.Cookie{
		Name:     "auth",
		Value:    "some-token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	req, _ := http.NewRequest("HEAD", fmt.Sprintf("%s/app/public/logout", s.URL), nil)
	req.AddCookie(c)
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

	cookie := test.GetCookieByName(resp.Cookies(), "auth")

	if cookie == nil {
		t.Fatal("Expect auth cookie, but no such cookie found")
	}

	if cookie.Expires.IsZero() {
		t.Fatal("Expected Expire to be set")
	}

	if !cookie.Expires.Before(time.Now()) {
		t.Fatal("Expected Expire to be before now, but it is not")
	}

	if cookie.Path != "/" {
		t.Fatalf("Expect path %s, got %s", "/", cookie.Path)
	}
}

func nominalAuthCheck(t *testing.T) {
	var payloadResp []byte
	userId := uuid.New()
	test.InsertUser(&model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	expTime := time.Now().Add(2 * time.Minute)

	tok := test.CreateToken(conf.JwtSecret, expTime, userId, t)

	c := &http.Cookie{
		Name:     "auth",
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expTime,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/session", s.URL), nil)
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		payloadResp, _ = ioutil.ReadAll(resp.Body)
		t.Fatalf("Expected 200 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}

	token := test.CheckAuthCookie(resp.Cookies(), conf, t)

	claims := token.Claims.(jwt.MapClaims)

	if claims["sub"] != userId.String() {
		t.Fatalf("Expected %s. Got %s", userId.String(), claims["sub"])
	}

	if claims["exp"].(float64) <= float64(expTime.Unix()) {
		t.Fatalf("Expect the  token to have been extended, but the exp is not greater than old one: %f vs %f", claims["exp"].(float64), float64(expTime.Unix()))
	}

	defer resp.Body.Close()
	payloadResp, _ = ioutil.ReadAll(resp.Body)

	var createdUser *response.User
	json.Unmarshal(payloadResp, &createdUser)

	if createdUser.Id != userId.String() {
		t.Fatalf("expect %s, got %s", userId.String(), createdUser.Id)
	}

	if createdUser.Login != "jerdct" {
		t.Fatalf("expect %s, got %s", "jerdct", createdUser.Login)
	}
}

func notTokenAuthCheck(t *testing.T) {
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/session", s.URL), nil)
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
}

func outDatedTokenAuthCheck(t *testing.T) {
	userId := uuid.New()
	test.InsertUser(&model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	expTime := time.Now().Add(-2 * time.Minute)

	tok := test.CreateToken(conf.JwtSecret, expTime, userId, t)

	c := &http.Cookie{
		Name:     "auth",
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expTime,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/session", s.URL), nil)
	req.AddCookie(c)
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

	cookie := test.GetCookieByName(resp.Cookies(), "auth")

	if cookie.Expires.After(time.Now()) {
		t.Fatal("Expect auth cookie to expire")
	}
}
