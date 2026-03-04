package conformvault

import (
	"context"
	"fmt"
	"io"
)

// AuditService handles audit log operations.
type AuditService struct {
	client *Client
}

// List returns audit log entries with optional filters.
func (s *AuditService) List(ctx context.Context, opts *AuditListOptions) ([]AuditEntry, error) {
	path := "/audit"
	if opts != nil {
		sep := "?"
		if opts.EventType != "" {
			path += sep + "event_type=" + opts.EventType
			sep = "&"
		}
		if opts.From != "" {
			path += sep + "from=" + opts.From
			sep = "&"
		}
		if opts.To != "" {
			path += sep + "to=" + opts.To
			sep = "&"
		}
		if opts.Page > 0 {
			path += sep + fmt.Sprintf("page=%d", opts.Page)
			sep = "&"
		}
		if opts.Limit > 0 {
			path += sep + fmt.Sprintf("limit=%d", opts.Limit)
		}
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[AuditEntry]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Search performs a full-text search on audit entries.
func (s *AuditService) Search(ctx context.Context, opts *AuditSearchOptions) ([]AuditEntry, error) {
	path := "/audit/search"
	if opts != nil {
		sep := "?"
		if opts.Query != "" {
			path += sep + "q=" + opts.Query
			sep = "&"
		}
		if opts.EventType != "" {
			path += sep + "event_type=" + opts.EventType
			sep = "&"
		}
		if opts.From != "" {
			path += sep + "from=" + opts.From
			sep = "&"
		}
		if opts.To != "" {
			path += sep + "to=" + opts.To
			sep = "&"
		}
		if opts.Page > 0 {
			path += sep + fmt.Sprintf("page=%d", opts.Page)
			sep = "&"
		}
		if opts.Limit > 0 {
			path += sep + fmt.Sprintf("limit=%d", opts.Limit)
		}
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[AuditEntry]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Export exports audit logs in the specified format.
func (s *AuditService) Export(ctx context.Context, opts *AuditExportOptions) (io.ReadCloser, error) {
	path := "/audit/export"
	if opts != nil {
		sep := "?"
		if opts.Format != "" {
			path += sep + "format=" + opts.Format
			sep = "&"
		}
		if opts.EventType != "" {
			path += sep + "event_type=" + opts.EventType
			sep = "&"
		}
		if opts.From != "" {
			path += sep + "from=" + opts.From
			sep = "&"
		}
		if opts.To != "" {
			path += sep + "to=" + opts.To
		}
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}

// GetStats returns aggregated audit statistics.
func (s *AuditService) GetStats(ctx context.Context) (*AuditStats, error) {
	req, err := s.client.newRequest(ctx, "GET", "/audit/stats", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[AuditStats]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetAnomalies returns detected audit anomalies.
func (s *AuditService) GetAnomalies(ctx context.Context) ([]AuditAnomaly, error) {
	req, err := s.client.newRequest(ctx, "GET", "/audit/anomalies", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[AuditAnomaly]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
