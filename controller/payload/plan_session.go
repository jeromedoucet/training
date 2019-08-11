package payload

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

type PlanSession struct {
	PlanId      string
	Day         *time.Time `json:"day"`
	Description string     `json:"description"`
}

func (p *PlanSession) ToModel() (*model.PlanSession, error) {
	var planId uuid.UUID
	var err error

	planId, err = uuid.Parse(p.PlanId)

	if err != nil {
		return nil, err
	}

	if len(p.Description) == 0 {
		return nil, errors.New("Expect a non empty description")
	}

	if p.Day == nil {
		return nil, errors.New("Expect a non nil day")
	}

	return &model.PlanSession{PlanId: planId, Day: *p.Day, Description: p.Description}, nil
}
