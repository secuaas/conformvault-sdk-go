package conformvault

import "context"

// JobsService manages background jobs.
type JobsService struct {
	client *Client
}

// Create creates a new background job.
func (s *JobsService) Create(ctx context.Context, request CreateJobRequest) (*Job, error) {
	req, err := s.client.newRequest(ctx, "POST", "/jobs", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Job]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all background jobs.
func (s *JobsService) List(ctx context.Context) ([]Job, error) {
	req, err := s.client.newRequest(ctx, "GET", "/jobs", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Job]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single background job by ID.
func (s *JobsService) Get(ctx context.Context, jobID string) (*Job, error) {
	req, err := s.client.newRequest(ctx, "GET", "/jobs/"+jobID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Job]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Cancel cancels a background job.
func (s *JobsService) Cancel(ctx context.Context, jobID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/jobs/"+jobID+"/cancel", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
