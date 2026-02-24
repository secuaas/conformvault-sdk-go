package conformvault

import (
	"context"
	"fmt"
)

// ShareLinksService handles share link operations.
type ShareLinksService struct {
	client *Client
}

// List returns all share links.
func (s *ShareLinksService) List(ctx context.Context, opts *ShareLinkListOptions) ([]ShareLink, error) {
	path := "/sharelinks"
	if opts != nil {
		sep := "?"
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

	var resp ListResponse[ShareLink]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Create creates a new share link.
func (s *ShareLinksService) Create(ctx context.Context, r CreateShareLinkRequest) (*ShareLink, error) {
	req, err := s.client.newRequest(ctx, "POST", "/sharelinks", r)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[ShareLink]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a share link by ID.
func (s *ShareLinksService) Delete(ctx context.Context, shareLinkID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/sharelinks/"+shareLinkID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
