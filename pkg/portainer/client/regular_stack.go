package client

import (
	"fmt"

	"github.com/portainer/client-api-go/v2/pkg/client/stacks"
	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
)

// ListRegularStacks lists all regular stacks using the raw API client
func (c *PortainerClient) ListRegularStacks() ([]*apimodels.PortainereeStack, error) {
	if c.rawCli == nil {
		return nil, fmt.Errorf("raw API client is not initialized")
	}

	params := stacks.NewStackListParams()
	okResp, noContentResp, err := c.rawCli.Stacks.StackList(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list regular stacks: %w", err)
	}

	// Handle success response
	if okResp != nil {
		return okResp.Payload, nil
	}

	// Handle no content response
	if noContentResp != nil {
		return []*apimodels.PortainereeStack{}, nil
	}

	return nil, fmt.Errorf("unexpected empty response from stack list")
}

// CreateRegularStack creates a new regular stack using the raw API client
func (c *PortainerClient) CreateRegularStack(name string, file string, environmentId int64) (int64, error) {
	if c.rawCli == nil {
		return 0, fmt.Errorf("raw API client is not initialized")
	}

	params := stacks.NewStackCreateDockerStandaloneStringParams().WithEndpointID(environmentId).WithBody(&apimodels.StacksComposeStackFromFileContentPayload{
		Name:             &name,
		StackFileContent: &file,
	})

	resp, err := c.rawCli.Stacks.StackCreateDockerStandaloneString(params, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create regular stack: %w", err)
	}

	// The method returns *StackCreateDockerStandaloneStringOK directly
	if resp != nil && resp.Payload != nil {
		return resp.Payload.ID, nil
	}

	return 0, fmt.Errorf("failed to get stack ID from create response")
}

// UpdateRegularStack updates an existing regular stack using the raw API client
func (c *PortainerClient) UpdateRegularStack(id int64, file string, environmentId int64) error {
	if c.rawCli == nil {
		return fmt.Errorf("raw API client is not initialized")
	}

	params := stacks.NewStackUpdateParams().WithID(id).WithEndpointID(environmentId).WithBody(&apimodels.StacksUpdateStackPayload{
		StackFileContent: file,
	})

	_, err := c.rawCli.Stacks.StackUpdate(params, nil)
	if err != nil {
		return fmt.Errorf("failed to update regular stack: %w", err)
	}

	return nil
}

// GetRegularStackFile gets the file content for a regular stack using the raw API client
func (c *PortainerClient) GetRegularStackFile(id int64) (string, error) {
	if c.rawCli == nil {
		return "", fmt.Errorf("raw API client is not initialized")
	}

	params := stacks.NewStackFileInspectParams().WithID(id)
	resp, err := c.rawCli.Stacks.StackFileInspect(params, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get regular stack file: %w", err)
	}

	// The method returns *StackFileInspectOK directly
	if resp != nil && resp.Payload != nil {
		return resp.Payload.StackFileContent, nil
	}

	return "", fmt.Errorf("stack file response payload is nil")
}