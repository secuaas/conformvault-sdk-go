package conformvault

import (
	"context"
	"fmt"
)

// BatchesService handles batch operation endpoints.
type BatchesService struct {
	client *Client
}

// --- Types ---

// BatchOperation represents a batch operation.
type BatchOperation struct {
	ID        string               `json:"id"`
	Status    string               `json:"status"`
	Type      string               `json:"type"`
	Total     int                  `json:"total"`
	Completed int                  `json:"completed"`
	Failed    int                  `json:"failed"`
	Items     []BatchOperationItem `json:"items,omitempty"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

// BatchOperationItem represents a single item within a batch operation.
type BatchOperationItem struct {
	ID       string  `json:"id"`
	Index    int     `json:"index"`
	Filename string  `json:"filename"`
	Size     int64   `json:"size"`
	MimeType string  `json:"mime_type"`
	Status   string  `json:"status"`
	FileID   *string `json:"file_id,omitempty"`
	Error    *string `json:"error,omitempty"`
}

// BatchOperationResponse wraps a single batch operation response.
type BatchOperationResponse struct {
	BatchOperation
}

// BatchListResponse wraps a paginated list of batch operations.
type BatchListResponse struct {
	Data  []BatchOperation `json:"data"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int              `json:"total"`
}

// CreateBatchRequest is the input for creating a batch operation.
type CreateBatchRequest struct {
	Type     string               `json:"type"`
	FolderID *string              `json:"folder_id,omitempty"`
	Items    []CreateBatchItemDef `json:"items"`
}

// CreateBatchItemDef defines a single item in a batch creation request.
type CreateBatchItemDef struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// --- Methods ---

// Create creates a new batch operation.
func (s *BatchesService) Create(ctx context.Context, r CreateBatchRequest) (*BatchOperationResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/batches", r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[BatchOperationResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of batch operations.
func (s *BatchesService) List(ctx context.Context, page, limit int) (*BatchListResponse, error) {
	path := "/batches"
	sep := "?"
	if page > 0 {
		path += fmt.Sprintf("%spage=%d", sep, page)
		sep = "&"
	}
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp BatchListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single batch operation by ID.
func (s *BatchesService) Get(ctx context.Context, id string) (*BatchOperationResponse, error) {
	req, err := s.client.newRequest(ctx, "GET", "/batches/"+id, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[BatchOperationResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Commit finalizes a batch operation and triggers processing.
func (s *BatchesService) Commit(ctx context.Context, id string) (*BatchOperationResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/batches/"+id+"/commit", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[BatchOperationResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Cancel cancels a batch operation.
func (s *BatchesService) Cancel(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, "POST", "/batches/"+id+"/cancel", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
