package conformvault

import (
	"context"
	"fmt"
)

// FoldersService handles folder operations.
type FoldersService struct {
	client *Client
}

// List returns all folders, optionally filtered.
func (s *FoldersService) List(ctx context.Context, opts *FolderListOptions) ([]Folder, error) {
	path := "/folders"
	if opts != nil {
		sep := "?"
		if opts.ParentID != nil {
			path += sep + "parent_id=" + *opts.ParentID
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

	var resp ListResponse[Folder]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single folder by ID.
func (s *FoldersService) Get(ctx context.Context, folderID string) (*Folder, error) {
	req, err := s.client.newRequest(ctx, "GET", "/folders/"+folderID, nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[Folder]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Create creates a new folder.
func (s *FoldersService) Create(ctx context.Context, r CreateFolderRequest) (*Folder, error) {
	req, err := s.client.newRequest(ctx, "POST", "/folders", r)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[Folder]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a folder by ID.
func (s *FoldersService) Delete(ctx context.Context, folderID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/folders/"+folderID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
