package conformvault

import (
	"context"
	"fmt"
	"time"
)

// SearchService handles unified search across files and folders.
type SearchService struct {
	client *Client
}

// SearchOptions are query parameters for search.
type SearchOptions struct {
	Query    string  `json:"q"`
	Types    string  `json:"types,omitempty"`    // "files", "folders", or "files,folders"
	FolderID *string `json:"folder_id,omitempty"`
	Page     int     `json:"page,omitempty"`
	PageSize int     `json:"page_size,omitempty"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"` // "file" or "folder"
	Name         string    `json:"name"`
	Path         string    `json:"path,omitempty"`
	Size         int64     `json:"size,omitempty"`
	ContentType  string    `json:"content_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// SearchResponse wraps the search response with pagination.
type SearchResponse struct {
	Data       []SearchResult `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination contains pagination metadata.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// Search performs a full-text search across files and folders.
func (s *SearchService) Search(ctx context.Context, opts SearchOptions) (*SearchResponse, error) {
	path := fmt.Sprintf("/search?q=%s", opts.Query)

	if opts.Types != "" {
		path += "&types=" + opts.Types
	}
	if opts.FolderID != nil {
		path += "&folder_id=" + *opts.FolderID
	}
	if opts.Page > 0 {
		path += fmt.Sprintf("&page=%d", opts.Page)
	}
	if opts.PageSize > 0 {
		path += fmt.Sprintf("&page_size=%d", opts.PageSize)
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
