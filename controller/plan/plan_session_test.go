package plan_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/controller"
	"github.com/jeromedoucet/training/test"
)

func TestPlanSessionCreationSuite(t *testing.T) {

	test.CleanDB(db)
	insertDataSet()

	t.Run("nominal plan session creation", nominalPlanSessionCreation)
	t.Run("plan session creation when not autheticated", nominalPlanSessionCreation)
	t.Run("missing field during plan session creation", missingFieldPlanSessionCreation)
	t.Run("no plan when tempting to create a plan session", PlanSession404)
	// todo, more than a Day, prefers two instant in time. begin, end.
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

}
