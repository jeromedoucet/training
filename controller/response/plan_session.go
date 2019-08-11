package response

import (
	"time"

	"github.com/jeromedoucet/training/model"
)

type PlanSession struct {
	Id          string    `json:"id"`
	PlanId      string    `json:"plan_id"`
	Day         time.Time `json:"day"`
	Description string    `json:"description"`
	Comments    string    `json:"comments"`
}

func FromPlanSessionModel(planSession *model.PlanSession) *PlanSession {
	return &PlanSession{Id: planSession.Id.String(), PlanId: planSession.PlanId.String(), Day: planSession.Day, Description: planSession.Description, Comments: planSession.Comments}
}
