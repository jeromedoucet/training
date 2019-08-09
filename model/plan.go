package model

import "github.com/google/uuid"

type Plan struct {
	Id        uuid.UUID
	Name      string
	CreatorId uuid.UUID
	TraineeId uuid.UUID
}
