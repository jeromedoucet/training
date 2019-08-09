package payload_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/model"
)

func TestToModel(t *testing.T) {
	var err error
	var res *model.Plan
	creatorId := uuid.New()
	traineeId := uuid.New()

	plan := &payload.Plan{}

	_, err = plan.ToModel()

	if err == nil {
		t.Fatal("Expect error got nil")
	}

	plan.Name = "training plan"
	_, err = plan.ToModel()

	if err == nil {
		t.Fatal("Expect error got nil")
	}

	plan.CreatorId = creatorId.String()

	_, err = plan.ToModel()

	if err == nil {
		t.Fatal("Expect error got nil")
	}

	plan.TraineeId = traineeId.String()

	res, err = plan.ToModel()

	if err != nil {
		t.Fatalf("Expect no error got %s", err.Error())
	}

	if res.Name != plan.Name {
		t.Fatalf("Expect %s to equal %s", res.Name, plan.Name)
	}

	if res.CreatorId.String() != plan.CreatorId {
		t.Fatalf("Expect %s to equal %s", res.CreatorId.String(), plan.CreatorId)
	}

	if res.TraineeId.String() != plan.TraineeId {
		t.Fatalf("Expect %s to equal %s", res.TraineeId.String(), plan.TraineeId)
	}
}
