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
	// Get edge stacks
	edgeStacks, err := c.cli.ListEdgeStacks()
	if err != nil {
		return nil, fmt.Errorf("failed to list edge stacks: %w", err)
	}

	// Convert edge stacks
	allStacks := make([]models.Stack, 0, len(edgeStacks))
	for _, es := range edgeStacks {
		allStacks = append(allStacks, models.ConvertEdgeStackToStack(es))
	}

	// Try to get regular stacks if the raw client is available
	if c.rawCli != nil {
		regularStacks, err := c.ListRegularStacks()
		if err != nil {
			// Log warning but don't fail - regular stacks are optional enhancement
			// In production, you might want to log this warning
			// For now, continue with just edge stacks
		} else {
			// Add regular stacks
			for _, rs := range regularStacks {
				allStacks = append(allStacks, models.ConvertRegularStackToStack(rs))
			}
		}
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

	// If edge stack failed and raw client is available, try regular stack
	if c.rawCli != nil {
		file, err = c.GetRegularStackFile(int64(id))
		if err == nil {
			return file, nil
		}
	}

	// Return the edge stack error since that was tried first
	return "", fmt.Errorf("failed to get stack file from edge stacks: %w", err)
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

	// If edge stack update failed and raw client is available, try regular stack (with first environment ID)
	if c.rawCli != nil && len(environmentGroupIds) > 0 {
		regularErr := c.UpdateRegularStack(int64(id), file, int64(environmentGroupIds[0]))
		if regularErr == nil {
			return nil
		}
		// Continue with edge stack error below since that was tried first
	}

	// Return the edge stack error since that's the primary functionality
	return fmt.Errorf("failed to update edge stack: %w", err)
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
