package conformvault

import (
	"context"
	"fmt"
)

// ScanReportsService handles file scan report operations.
type ScanReportsService struct {
	client *Client
}

// GetReport returns the scan report for a specific file.
func (s *ScanReportsService) GetReport(ctx context.Context, fileID string) (*FileScanReport, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/scan-report", nil)
	if err != nil {
		return nil, err
	}
	var report FileScanReport
	if err := s.client.do(req, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// List returns scan reports with optional pagination.
func (s *ScanReportsService) List(ctx context.Context, opts *ScanReportListOptions) (*ScanReportListResponse, error) {
	path := "/scan-reports"
	if opts != nil {
		sep := "?"
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
		return nil, err
	}
	var resp ScanReportListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSummary returns aggregate scan statistics for the organization.
func (s *ScanReportsService) GetSummary(ctx context.Context) (*FileScanSummary, error) {
	req, err := s.client.newRequest(ctx, "GET", "/scan-reports/summary", nil)
	if err != nil {
		return nil, err
	}
	var summary FileScanSummary
	if err := s.client.do(req, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}
