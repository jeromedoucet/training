package payload

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

type PlanSession struct {
	PlanId      string
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	Description string    `json:"description"`
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

	if p.From.IsZero() {
		return nil, errors.New("Expect a non nil day")
	}

	if p.To.IsZero() {
		return nil, errors.New("Expect a non nil day")
	}

	return &model.PlanSession{PlanId: planId, From: p.From, To: p.To, Description: p.Description}, nil
}
