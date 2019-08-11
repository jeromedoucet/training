package model

import (
	"time"

	"github.com/google/uuid"
)

type PlanSession struct {
	Id          uuid.UUID
	PlanId      uuid.UUID
	Day         time.Time
	Description string
	Comments    string
}
