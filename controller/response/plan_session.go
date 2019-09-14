package response

import (
	"time"

	"github.com/jeromedoucet/training/model"
)

type PlanSession struct {
	Id          string    `json:"id"`
	PlanId      string    `json:"plan_id"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	Description string    `json:"description"`
	Comments    string    `json:"comments"`
}

func FromPlanSessionModel(planSession *model.PlanSession) *PlanSession {
	return &PlanSession{Id: planSession.Id.String(), PlanId: planSession.PlanId.String(), From: planSession.From, To: planSession.To, Description: planSession.Description, Comments: planSession.Comments}
}
