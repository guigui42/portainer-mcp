package client

import (
	"fmt"

	"github.com/portainer/portainer-mcp/pkg/portainer/models"
	"github.com/portainer/portainer-mcp/pkg/portainer/utils"
)

// GetStacks retrieves all stacks from the Portainer server.
// This includes both Edge Stacks and regular Stacks.
//
// Returns:
//   - A slice of Stack objects
//   - An error if the operation fails
func (c *PortainerClient) GetStacks() ([]models.Stack, error) {
	// Get both edge stacks and regular stacks
	edgeStacks, err := c.cli.ListEdgeStacks()
	if err != nil {
		return nil, fmt.Errorf("failed to list edge stacks: %w", err)
	}

	regularStacks, err := c.ListRegularStacks()
	if err != nil {
		return nil, fmt.Errorf("failed to list regular stacks: %w", err)
	}

	// Convert and combine the results
	allStacks := make([]models.Stack, 0, len(edgeStacks)+len(regularStacks))

	// Add edge stacks
	for _, es := range edgeStacks {
		allStacks = append(allStacks, models.ConvertEdgeStackToStack(es))
	}

	// Add regular stacks
	for _, rs := range regularStacks {
		allStacks = append(allStacks, models.ConvertRegularStackToStack(rs))
	}

	return allStacks, nil
}

// GetStackFile retrieves the file content of a stack from the Portainer server.
// This method attempts to get the file from both Edge Stacks and regular Stacks.
//
// Parameters:
//   - id: The ID of the stack to retrieve
//
// Returns:
//   - The file content of the stack (Compose file)
//   - An error if the operation fails
func (c *PortainerClient) GetStackFileContent(id int) (string, error) {
	// Try to get from edge stack first
	file, err := c.cli.GetEdgeStackFile(int64(id))
	if err == nil {
		return file, nil
	}

	// If edge stack failed, try regular stack
	file, err = c.GetRegularStackFile(int64(id))
	if err != nil {
		return "", fmt.Errorf("failed to get stack file from both edge and regular stacks: %w", err)
	}

	return file, nil
}

// CreateStackWrapper creates a new stack on the Portainer server.
// This function creates an Edge Stack if environmentGroupIds are provided,
// otherwise creates a regular Stack.
//
// Parameters:
//   - name: The name of the stack
//   - file: The file content of the stack (Compose file)
//   - environmentGroupIds: A slice of environment group IDs (for edge stacks) or environments (for regular stacks)
//
// Returns:
//   - The ID of the created stack
//   - An error if the operation fails
func (c *PortainerClient) CreateStackWrapper(name, file string, environmentGroupIds []int) (int, error) {
	// For now, maintain backward compatibility by creating Edge Stacks
	// TODO: Add parameter or logic to determine whether to create regular or edge stack
	id, err := c.cli.CreateEdgeStack(name, file, utils.IntToInt64Slice(environmentGroupIds))
	if err != nil {
		return 0, fmt.Errorf("failed to create edge stack: %w", err)
	}

	return int(id), nil
}

// UpdateStackWrapper updates an existing stack on the Portainer server.
// This function updates either an Edge Stack or regular Stack based on the stack ID.
//
// Parameters:
//   - id: The ID of the stack to update
//   - file: The file content of the stack (Compose file)
//   - environmentGroupIds: A slice of environment group IDs to include in the stack
//
// Returns:
//   - An error if the operation fails
func (c *PortainerClient) UpdateStackWrapper(id int, file string, environmentGroupIds []int) error {
	// Try to update as edge stack first
	err := c.cli.UpdateEdgeStack(int64(id), file, utils.IntToInt64Slice(environmentGroupIds))
	if err == nil {
		return nil
	}

	// If edge stack update failed, try regular stack (with first environment ID)
	if len(environmentGroupIds) > 0 {
		err = c.UpdateRegularStack(int64(id), file, int64(environmentGroupIds[0]))
		if err != nil {
			return fmt.Errorf("failed to update stack as both edge and regular stack: %w", err)
		}
		return nil
	}

	return fmt.Errorf("failed to update edge stack and no environment ID provided for regular stack: %w", err)
}

// GetStackFile is an alias for GetStackFileContent to match the expected interface
func (c *PortainerClient) GetStackFile(id int) (string, error) {
	return c.GetStackFileContent(id)
}

// CreateStack is an alias for CreateStackWrapper to match the expected interface
func (c *PortainerClient) CreateStack(name, file string, environmentGroupIds []int) (int, error) {
	return c.CreateStackWrapper(name, file, environmentGroupIds)
}

// UpdateStack is an alias for UpdateStackWrapper to match the expected interface
func (c *PortainerClient) UpdateStack(id int, file string, environmentGroupIds []int) error {
	return c.UpdateStackWrapper(id, file, environmentGroupIds)
}
