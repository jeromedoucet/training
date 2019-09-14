package model

import (
	"time"

	"github.com/google/uuid"
)

type PlanSession struct {
	Id          uuid.UUID
	PlanId      uuid.UUID
	From        time.Time
	To          time.Time
	Description string
	Comments    string
}
