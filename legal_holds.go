package conformvault

import "context"

// LegalHoldsService manages legal holds.
type LegalHoldsService struct {
	client *Client
}

// Create creates a new legal hold.
func (s *LegalHoldsService) Create(ctx context.Context, request CreateLegalHoldRequest) (*LegalHold, error) {
	req, err := s.client.newRequest(ctx, "POST", "/legal-holds", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[LegalHold]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all legal holds.
func (s *LegalHoldsService) List(ctx context.Context) ([]LegalHold, error) {
	req, err := s.client.newRequest(ctx, "GET", "/legal-holds", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[LegalHold]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single legal hold by ID.
func (s *LegalHoldsService) Get(ctx context.Context, holdID string) (*LegalHold, error) {
	req, err := s.client.newRequest(ctx, "GET", "/legal-holds/"+holdID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[LegalHold]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Release releases a legal hold.
func (s *LegalHoldsService) Release(ctx context.Context, holdID string) (*LegalHold, error) {
	req, err := s.client.newRequest(ctx, "POST", "/legal-holds/"+holdID+"/release", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[LegalHold]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddFiles adds files to a legal hold.
func (s *LegalHoldsService) AddFiles(ctx context.Context, holdID string, request AddLegalHoldFilesRequest) ([]LegalHoldFile, error) {
	req, err := s.client.newRequest(ctx, "POST", "/legal-holds/"+holdID+"/files", request)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[LegalHoldFile]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// RemoveFile removes a file from a legal hold.
func (s *LegalHoldsService) RemoveFile(ctx context.Context, holdID, fileID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/legal-holds/"+holdID+"/files/"+fileID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
