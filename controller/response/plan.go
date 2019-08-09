package response

import (
	"github.com/jeromedoucet/training/model"
)

// Plan is the structure used for
// the serialization of training entity
type Plan struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatorId string `json:"creator_id"`
	TraineeId string `json:"trainee_id"`
}

func FromPlanModel(plan *model.Plan) *Plan {
	return &Plan{Id: plan.Id.String(), CreatorId: plan.CreatorId.String(), TraineeId: plan.TraineeId.String(), Name: plan.Name}
}
