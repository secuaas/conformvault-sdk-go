package conformvault

import (
	"context"
	"fmt"
)

// TrashService handles trash / recycle bin operations.
type TrashService struct {
	client *Client
}

// TrashListOptions are query parameters for listing trashed files.
type TrashListOptions struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// TrashResponse wraps the trash list response with pagination.
type TrashResponse struct {
	Data       []File     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// EmptyTrashResponse is the response from emptying the trash.
type EmptyTrashResponse struct {
	Message      string `json:"message"`
	FilesDeleted int    `json:"files_deleted"`
}

// List returns all files in the trash.
func (s *TrashService) List(ctx context.Context, opts *TrashListOptions) (*TrashResponse, error) {
	path := "/trash"
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

	var resp TrashResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Restore restores a file from the trash.
func (s *TrashService) Restore(ctx context.Context, fileID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/trash/"+fileID+"/restore", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Delete permanently deletes a file from the trash.
func (s *TrashService) Delete(ctx context.Context, fileID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/trash/"+fileID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Empty permanently deletes all files in the trash.
func (s *TrashService) Empty(ctx context.Context) (*EmptyTrashResponse, error) {
	req, err := s.client.newRequest(ctx, "DELETE", "/trash", nil)
	if err != nil {
		return nil, err
	}

	var resp EmptyTrashResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
