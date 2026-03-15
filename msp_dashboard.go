package conformvault

import (
	"context"
	"fmt"
	"net/url"
)

// MSPDashboardService provides MSP dashboard and client metrics.
type MSPDashboardService struct {
	client *Client
}

// mspClientListResponse wraps paginated MSP client metrics with a total count.
type mspClientListResponse struct {
	Data  []MSPClientMetrics `json:"data"`
	Total int                `json:"total"`
}

// GetDashboard returns the MSP dashboard summary.
func (s *MSPDashboardService) GetDashboard(ctx context.Context) (*MSPDashboardSummary, error) {
	req, err := s.client.newRequest(ctx, "GET", "/msp/dashboard", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[MSPDashboardSummary]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListClients returns MSP client metrics with optional search and pagination.
func (s *MSPDashboardService) ListClients(ctx context.Context, opts *MSPClientListOptions) ([]MSPClientMetrics, int, error) {
	path := "/msp/clients"
	if opts != nil {
		sep := "?"
		if opts.Search != "" {
			path += sep + "search=" + url.QueryEscape(opts.Search)
			sep = "&"
		}
		if opts.Limit > 0 {
			path += sep + fmt.Sprintf("limit=%d", opts.Limit)
			sep = "&"
		}
		if opts.Offset > 0 {
			path += sep + fmt.Sprintf("offset=%d", opts.Offset)
		}
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, 0, err
	}
	var resp mspClientListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Data, resp.Total, nil
}

// GetClientUsage returns detailed usage metrics for a specific MSP client.
func (s *MSPDashboardService) GetClientUsage(ctx context.Context, clientID string) (*MSPClientUsage, error) {
	req, err := s.client.newRequest(ctx, "GET", "/msp/clients/"+clientID+"/usage", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[MSPClientUsage]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
