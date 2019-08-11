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
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/model"
	"github.com/jeromedoucet/training/test"
)

func TestPlanSuite(t *testing.T) {

	test.CleanDB(db)
	t.Run("nominal training creation", nominalPlanCreation)

	test.CleanDB(db)
	t.Run("missing field training creation", missingFieldPlanCreation)

	test.CleanDB(db)
	t.Run("no auth training creation", notAuthenticatedPlanCreation)
}

func nominalPlanCreation(t *testing.T) {
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

	payload := struct {
		Name      string `json:"name"`
		CreatorId string `json:"creator_id"`
		TraineeId string `json:"trainee_id"`
	}{
		"Training plan",
		userId.String(),
		userId.String(),
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan", s.URL), bytes.NewBuffer(body))
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

	payloadResp, _ = ioutil.ReadAll(resp.Body)
	var createdPlan *response.Plan
	json.Unmarshal(payloadResp, &createdPlan)

	if len(createdPlan.Id) == 0 {
		t.Fatal("expect non empty id")
	}

	if createdPlan.Name != payload.Name {
		t.Fatalf("Expect %s, got %s", payload.Name, createdPlan.Name)
	}

	if createdPlan.CreatorId != payload.CreatorId {
		t.Fatalf("Expect %s, got %s", payload.CreatorId, createdPlan.CreatorId)
	}

	if createdPlan.TraineeId != payload.TraineeId {
		t.Fatalf("Expect %s, got %s", payload.TraineeId, createdPlan.TraineeId)
	}
}

func missingFieldPlanCreation(t *testing.T) {
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

	payload := struct {
		Name      string `json:"name"`
		TraineeId string `json:"trainee_id"`
	}{
		"Training plan",
		userId.String(),
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan", s.URL), bytes.NewBuffer(body))
	req.AddCookie(c)
	client := &http.Client{}

	// when
	resp, err := client.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expected to have no error, but got %s", err.Error())
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 return code. Got %d ", resp.StatusCode)
	}
}

func notAuthenticatedPlanCreation(t *testing.T) {
	userId := uuid.New()
	test.InsertUser(&model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}, db)
	s := httptest.NewServer(controller.InitRoutes(conf))
	defer s.Close()

	payload := struct {
		Name      string `json:"name"`
		CreatorId string `json:"creator_id"`
		TraineeId string `json:"trainee_id"`
	}{
		"Training plan",
		userId.String(),
		userId.String(),
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/app/public/plan", s.URL), bytes.NewBuffer(body))
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
