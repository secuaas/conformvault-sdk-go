package conformvault

import "context"

// ImportsService manages cloud storage import connections and jobs.
type ImportsService struct {
	client *Client
}

// CreateConnection creates a new import connection.
func (s *ImportsService) CreateConnection(ctx context.Context, request CreateImportConnectionRequest) (*ImportConnection, error) {
	req, err := s.client.newRequest(ctx, "POST", "/import/connections", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ImportConnection]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListConnections returns all import connections.
func (s *ImportsService) ListConnections(ctx context.Context) ([]ImportConnection, error) {
	req, err := s.client.newRequest(ctx, "GET", "/import/connections", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[ImportConnection]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeleteConnection deletes an import connection by ID.
func (s *ImportsService) DeleteConnection(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/import/connections/"+id, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// StartImport starts an import job from a connection.
func (s *ImportsService) StartImport(ctx context.Context, connectionID string, request StartImportRequest) (*ImportJob, error) {
	req, err := s.client.newRequest(ctx, "POST", "/import/connections/"+connectionID+"/start", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ImportJob]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListJobs returns all import jobs.
func (s *ImportsService) ListJobs(ctx context.Context) ([]ImportJob, error) {
	req, err := s.client.newRequest(ctx, "GET", "/import/jobs", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[ImportJob]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetJob returns a single import job by ID.
func (s *ImportsService) GetJob(ctx context.Context, jobID string) (*ImportJob, error) {
	req, err := s.client.newRequest(ctx, "GET", "/import/jobs/"+jobID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ImportJob]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CancelJob cancels a running import job.
func (s *ImportsService) CancelJob(ctx context.Context, jobID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/import/jobs/"+jobID+"/cancel", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
