package payload_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/model"
)

func TestPlanSessionToModel(t *testing.T) {
	var err error
	var actualPlanSession *model.PlanSession

	planSession := &payload.PlanSession{}

	actualPlanSession, err = planSession.ToModel()

	if err == nil {
		t.Fatal("Expect non nil err")
	}

	now := time.Now()
	planSession.From = now

	actualPlanSession, err = planSession.ToModel()

	if err == nil {
		t.Fatal("Expect non nil err")
	}

	planSession.To = now.Add(time.Hour)

	actualPlanSession, err = planSession.ToModel()

	if err == nil {
		t.Fatal("Expect non nil err")
	}

	planSession.Description = "some description"

	actualPlanSession, err = planSession.ToModel()

	if err == nil {
		t.Fatal("Expect non nil err")
	}

	planSession.PlanId = uuid.New().String()
	actualPlanSession, err = planSession.ToModel()

	if err != nil {
		t.Fatalf("Expect non nil error, got %s", err.Error())
	}

	if actualPlanSession.Description != planSession.Description {
		t.Fatalf("Expect %s, got %s", planSession.Description, actualPlanSession.Description)
	}

	if !actualPlanSession.From.Equal(planSession.From) {
		t.Fatalf("Expect %v, got %v", planSession.From, actualPlanSession.From)
	}

	if !actualPlanSession.To.Equal(planSession.To) {
		t.Fatalf("Expect %v, got %v", planSession.To, actualPlanSession.To)
	}
}
