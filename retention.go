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

// CreateException creates an exception to a retention policy.
func (s *RetentionService) CreateException(ctx context.Context, policyID string, request CreateRetentionExceptionRequest) (*RetentionException, error) {
	req, err := s.client.newRequest(ctx, "POST", "/retention-policies/"+policyID+"/exceptions", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[RetentionException]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListExceptions returns all exceptions for a retention policy.
func (s *RetentionService) ListExceptions(ctx context.Context, policyID string) ([]RetentionException, error) {
	req, err := s.client.newRequest(ctx, "GET", "/retention-policies/"+policyID+"/exceptions", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[RetentionException]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeleteException deletes a retention exception by ID.
func (s *RetentionService) DeleteException(ctx context.Context, exceptionID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/retention-policies/exceptions/"+exceptionID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// ListPendingApprovals returns all pending retention approvals.
func (s *RetentionService) ListPendingApprovals(ctx context.Context) ([]RetentionApproval, error) {
	req, err := s.client.newRequest(ctx, "GET", "/retention-approvals", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[RetentionApproval]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DecideApproval approves or rejects a retention approval.
func (s *RetentionService) DecideApproval(ctx context.Context, approvalID string, request DecideApprovalRequest) error {
	req, err := s.client.newRequest(ctx, "POST", "/retention-approvals/"+approvalID+"/decide", request)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
