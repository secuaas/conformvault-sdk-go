package conformvault

import "context"

// RetentionService manages retention policies.
type RetentionService struct {
	client *Client
}

// Create creates a new retention policy.
func (s *RetentionService) Create(ctx context.Context, request CreateRetentionPolicyRequest) (*RetentionPolicy, error) {
	req, err := s.client.newRequest(ctx, "POST", "/retention-policies", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[RetentionPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all retention policies.
func (s *RetentionService) List(ctx context.Context) ([]RetentionPolicy, error) {
	req, err := s.client.newRequest(ctx, "GET", "/retention-policies", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[RetentionPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single retention policy by ID.
func (s *RetentionService) Get(ctx context.Context, policyID string) (*RetentionPolicy, error) {
	req, err := s.client.newRequest(ctx, "GET", "/retention-policies/"+policyID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[RetentionPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing retention policy.
func (s *RetentionService) Update(ctx context.Context, policyID string, request UpdateRetentionPolicyRequest) (*RetentionPolicy, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/retention-policies/"+policyID, request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[RetentionPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a retention policy by ID.
func (s *RetentionService) Delete(ctx context.Context, policyID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/retention-policies/"+policyID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
