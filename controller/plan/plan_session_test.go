package plan_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/controller"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/test"
)

func TestPlanSessionCreationSuite(t *testing.T) {

	test.CleanDB(db)
	insertDataSet()

	t.Run("nominal plan session creation", nominalPlanSessionCreation)
	t.Run("plan session creation when not autheticated", nominalPlanSessionCreation)
	t.Run("missing field during plan session creation", missingFieldPlanSessionCreation)
	t.Run("no plan when tempting to create a plan session", PlanSession404)
}

func nominalPlanSessionCreation(t *testing.T) {
	var payloadResp []byte
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

	payload := struct {
		From        time.Time `json:"from"`
		To          time.Time `json:"to"`
		Description string    `json:"description"`
	}{
		time.Now(),
		time.Now().Add(time.Hour * 2),
		"10 X 1' VMA",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, plan.Id.String()), bytes.NewBuffer(body))
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}
}

func PlanSessionNotAuthenticated(t *testing.T) {
	var payloadResp []byte
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		From        time.Time `json:"from"`
		To          time.Time `json:"to"`
		Description string    `json:"description"`
	}{
		time.Now(),
		time.Now().Add(time.Hour * 2),
		"10 X 1' VMA",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, plan.Id.String()), bytes.NewBuffer(body))
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}
}

func missingFieldPlanSessionCreation(t *testing.T) {
	var payloadResp []byte
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

	payload := struct {
		Description string `json:"description"`
	}{
		"10 X 1' VMA",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, plan.Id.String()), bytes.NewBuffer(body))
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}
}

func PlanSession404(t *testing.T) {
	var payloadResp []byte
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

	payload := struct {
		From        time.Time `json:"from"`
		To          time.Time `json:"to"`
		Description string    `json:"description"`
	}{
		time.Now(),
		time.Now().Add(time.Hour * 2),
		"10 X 1' VMA",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, uuid.New().String()), bytes.NewBuffer(body))
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}
}

func TestPlanSessionReadSuite(t *testing.T) {

	test.CleanDB(db)
	insertDataSet()

	t.Run("find to plan sessions with bounds", nominalPlanSessionFindWithBound)
	t.Run("find to plan sessions without bounds", nominalPlanSessionFindWithoutBound)
	t.Run("find to plan sessions unknown plan", nominalPlanSessionFindUnknownPlan)

	// todo 404
	// todo 401
	// todo 403

}

func nominalPlanSessionFindWithBound(t *testing.T) {
	var payloadResp []byte
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

	from := time.Date(2018, time.September, 30, 0, 0, 0, 0, time.UTC)
	to := time.Date(2018, time.October, 4, 0, 0, 0, 0, time.UTC)

	encFrom, _ := from.MarshalText()
	encTo, _ := to.MarshalText()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/plan/%s/sessions?from=%s&to=%s", s.URL, plan.Id.String(), base64.StdEncoding.EncodeToString(encFrom), base64.StdEncoding.EncodeToString(encTo)), nil)
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}

	var plans []*response.PlanSession

	json.Unmarshal(payloadResp, &plans)

	if len(plans) != 2 {
		t.Fatalf("Expect %d, got %d", 2, len(plans))
	}

	if plans[0].Id != firstSession.Id.String() {
		t.Fatal("Expect to get the first session, but get another")
	}

	if !plans[0].From.Equal(firstSession.From) {
		t.Fatalf("Expect %s, got %s", firstSession.From.String(), plans[0].From.String())
	}

	if !plans[0].To.Equal(firstSession.To) {
		t.Fatalf("Expect %s, got %s", firstSession.To.String(), plans[0].To.String())
	}

	if plans[0].Description != firstSession.Description {
		t.Fatalf("Expect %s, got %s", firstSession.Description, plans[0].Description)
	}

	if plans[0].Comments != firstSession.Comments {
		t.Fatalf("Expect %s, got %s", firstSession.Comments, plans[0].Comments)
	}

	if plans[1].Id != secondSession.Id.String() {
		t.Fatal("Expect to get the second session, but get another")
	}

	if !plans[1].From.Equal(secondSession.From) {
		t.Fatalf("Expect %s, got %s", secondSession.From.String(), plans[1].From.String())
	}

	if !plans[1].To.Equal(secondSession.To) {
		t.Fatalf("Expect %s, got %s", secondSession.To.String(), plans[1].To.String())
	}

	if plans[1].Description != secondSession.Description {
		t.Fatalf("Expect %s, got %s", secondSession.Description, plans[1].Description)
	}

	if plans[1].Comments != secondSession.Comments {
		t.Fatalf("Expect %s, got %s", secondSession.Comments, plans[1].Comments)
	}
}

func nominalPlanSessionFindWithoutBound(t *testing.T) {
	var payloadResp []byte
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

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, plan.Id.String()), nil)
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}

	var plans []*response.PlanSession

	json.Unmarshal(payloadResp, &plans)

	if len(plans) != 3 {
		t.Fatalf("Expect %d, got %d", 3, len(plans))
	}

	if plans[0].Id != firstSession.Id.String() {
		t.Fatal("Expect to get the first session, but get another")
	}

	if plans[1].Id != secondSession.Id.String() {
		t.Fatal("Expect to get the second session, but get another")
	}

	if plans[2].Id != thirdSession.Id.String() {
		t.Fatal("Expect to get the second session, but get another")
	}
}

func nominalPlanSessionFindUnknownPlan(t *testing.T) {
	var payloadResp []byte
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

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/app/public/plan/%s/sessions", s.URL, uuid.New().String()), nil)
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 return code. Got %d with body %s", resp.StatusCode, string(payloadResp))
	}

}
