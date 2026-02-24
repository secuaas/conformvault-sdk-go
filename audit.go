package conformvault

import (
	"context"
	"fmt"
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
