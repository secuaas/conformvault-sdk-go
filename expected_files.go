package conformvault

import "context"

// ExpectedFilesService manages expected files on share links.
type ExpectedFilesService struct {
	client *Client
}

// Create creates an expected file on a share link.
func (s *ExpectedFilesService) Create(ctx context.Context, shareLinkID string, request CreateExpectedFileRequest) (*ExpectedFile, error) {
	req, err := s.client.newRequest(ctx, "POST", "/sharelinks/"+shareLinkID+"/expected-files", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ExpectedFile]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all expected files for a share link.
func (s *ExpectedFilesService) List(ctx context.Context, shareLinkID string) ([]ExpectedFile, error) {
	req, err := s.client.newRequest(ctx, "GET", "/sharelinks/"+shareLinkID+"/expected-files", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[ExpectedFile]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Update updates an expected file.
func (s *ExpectedFilesService) Update(ctx context.Context, shareLinkID, efID string, request UpdateExpectedFileRequest) (*ExpectedFile, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/sharelinks/"+shareLinkID+"/expected-files/"+efID, request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ExpectedFile]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes an expected file.
func (s *ExpectedFilesService) Delete(ctx context.Context, shareLinkID, efID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/sharelinks/"+shareLinkID+"/expected-files/"+efID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetProgress returns the fulfillment progress for a share link's expected files.
func (s *ExpectedFilesService) GetProgress(ctx context.Context, shareLinkID string) (*ExpectedFileProgress, error) {
	req, err := s.client.newRequest(ctx, "GET", "/sharelinks/"+shareLinkID+"/expected-files/progress", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ExpectedFileProgress]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
