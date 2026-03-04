package conformvault

import "context"

// QuotaService provides quota information.
type QuotaService struct {
	client *Client
}

// Get returns the current quota usage.
func (s *QuotaService) Get(ctx context.Context) (*QuotaInfo, error) {
	req, err := s.client.newRequest(ctx, "GET", "/quota", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[QuotaInfo]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RateLimitService provides rate limit information.
type RateLimitService struct {
	client *Client
}

// Get returns the current rate limit status.
func (s *RateLimitService) Get(ctx context.Context) (*RateLimitInfo, error) {
	req, err := s.client.newRequest(ctx, "GET", "/rate-limit", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[RateLimitInfo]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
