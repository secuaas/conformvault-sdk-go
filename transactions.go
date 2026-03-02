package conformvault

import (
	"context"
	"fmt"
)

// TransactionsService handles transaction folder operations.
type TransactionsService struct {
	client *Client
}

// --- Types ---

// TransactionFolder represents a transaction folder.
type TransactionFolder struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Description    *string                 `json:"description,omitempty"`
	Status         string                  `json:"status"`
	DueDate        *string                 `json:"due_date,omitempty"`
	Progress       *TransactionProgress    `json:"progress,omitempty"`
	Items          []TransactionFolderItem `json:"items,omitempty"`
	CreatedAt      string                  `json:"created_at"`
	UpdatedAt      string                  `json:"updated_at"`
}

// TransactionFolderItem represents a single item in a transaction folder.
type TransactionFolderItem struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	Label         string  `json:"label"`
	Description   *string `json:"description,omitempty"`
	Required      bool    `json:"required"`
	Status        string  `json:"status"`
	FileID        *string `json:"file_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// TransactionProgress contains completion statistics for a transaction folder.
type TransactionProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Pending   int `json:"pending"`
}

// TransactionFolderResponse wraps a single transaction folder response.
type TransactionFolderResponse struct {
	TransactionFolder
}

// TransactionListResponse wraps a paginated list of transaction folders.
type TransactionListResponse struct {
	Data  []TransactionFolder `json:"data"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Total int                 `json:"total"`
}

// CreateTransactionRequest is the input for creating a transaction folder.
type CreateTransactionRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

// UpdateTransactionRequest is the input for updating a transaction folder.
type UpdateTransactionRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

// CreateTransactionItemRequest is the input for adding an item to a transaction folder.
type CreateTransactionItemRequest struct {
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Required    *bool   `json:"required,omitempty"`
}

// UpdateTransactionItemRequest is the input for updating a transaction folder item.
type UpdateTransactionItemRequest struct {
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	Required    *bool   `json:"required,omitempty"`
	Status      *string `json:"status,omitempty"`
	FileID      *string `json:"file_id,omitempty"`
}

// --- Methods ---

// Create creates a new transaction folder.
func (s *TransactionsService) Create(ctx context.Context, r CreateTransactionRequest) (*TransactionFolderResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/transactions", r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[TransactionFolderResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of transaction folders.
func (s *TransactionsService) List(ctx context.Context, page, limit int) (*TransactionListResponse, error) {
	path := "/transactions"
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
	var resp TransactionListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single transaction folder by ID.
func (s *TransactionsService) Get(ctx context.Context, id string) (*TransactionFolderResponse, error) {
	req, err := s.client.newRequest(ctx, "GET", "/transactions/"+id, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[TransactionFolderResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates a transaction folder.
func (s *TransactionsService) Update(ctx context.Context, id string, r UpdateTransactionRequest) (*TransactionFolderResponse, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/transactions/"+id, r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[TransactionFolderResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a transaction folder.
func (s *TransactionsService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/transactions/"+id, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// AddItem adds an item to a transaction folder.
func (s *TransactionsService) AddItem(ctx context.Context, txID string, r CreateTransactionItemRequest) (*TransactionFolderItem, error) {
	req, err := s.client.newRequest(ctx, "POST", "/transactions/"+txID+"/items", r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[TransactionFolderItem]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateItem updates an item in a transaction folder.
func (s *TransactionsService) UpdateItem(ctx context.Context, txID, itemID string, r UpdateTransactionItemRequest) (*TransactionFolderItem, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/transactions/"+txID+"/items/"+itemID, r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[TransactionFolderItem]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteItem removes an item from a transaction folder.
func (s *TransactionsService) DeleteItem(ctx context.Context, txID, itemID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/transactions/"+txID+"/items/"+itemID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
