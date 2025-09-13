package models

import (
	"time"

	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
	"github.com/portainer/portainer-mcp/pkg/portainer/utils"
)

type Stack struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	CreatedAt           string `json:"created_at"`
	EnvironmentGroupIds []int  `json:"group_ids"`
	StackType           string `json:"stack_type"` // "edge" or "regular"
	EnvironmentId       *int   `json:"environment_id,omitempty"` // Only for regular stacks
}

func ConvertEdgeStackToStack(rawEdgeStack *apimodels.PortainereeEdgeStack) Stack {
	createdAt := time.Unix(rawEdgeStack.CreationDate, 0).Format(time.RFC3339)

	return Stack{
		ID:                  int(rawEdgeStack.ID),
		Name:                rawEdgeStack.Name,
		CreatedAt:           createdAt,
		EnvironmentGroupIds: utils.Int64ToIntSlice(rawEdgeStack.EdgeGroups),
		StackType:           "edge",
	}
}

func ConvertRegularStackToStack(rawStack *apimodels.PortainereeStack) Stack {
	createdAt := time.Unix(rawStack.CreationDate, 0).Format(time.RFC3339)
	envId := int(rawStack.EndpointID)

	return Stack{
		ID:                  int(rawStack.ID),
		Name:                rawStack.Name,
		CreatedAt:           createdAt,
		EnvironmentGroupIds: []int{}, // Regular stacks don't use environment groups
		StackType:           "regular",
		EnvironmentId:       &envId,
	}
}
