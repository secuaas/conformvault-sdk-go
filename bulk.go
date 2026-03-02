package conformvault

import (
	"context"
	"io"
)

// BulkService handles bulk file operations.
type BulkService struct {
	client *Client
}

// BulkDeleteRequest is the input for bulk deleting files.
type BulkDeleteRequest struct {
	FileIDs []string `json:"file_ids"`
}

// BulkMoveRequest is the input for bulk moving files.
type BulkMoveRequest struct {
	FileIDs        []string `json:"file_ids"`
	TargetFolderID string   `json:"target_folder_id"`
}

// BulkDownloadRequest is the input for bulk downloading files as ZIP.
type BulkDownloadRequest struct {
	FileIDs []string `json:"file_ids"`
}

// BulkResult is the response from a bulk operation.
type BulkResult struct {
	Processed int           `json:"processed"`
	Succeeded int           `json:"succeeded"`
	Failed    int           `json:"failed"`
	Errors    []BulkError   `json:"errors,omitempty"`
}

// BulkError describes a single failure in a bulk operation.
type BulkError struct {
	FileID string `json:"file_id"`
	Error  string `json:"error"`
}

// Delete soft-deletes multiple files at once.
func (s *BulkService) Delete(ctx context.Context, fileIDs []string) (*BulkResult, error) {
	req, err := s.client.newRequest(ctx, "POST", "/files/bulk-delete", BulkDeleteRequest{FileIDs: fileIDs})
	if err != nil {
		return nil, err
	}

	var resp BulkResult
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Move moves multiple files to a target folder.
func (s *BulkService) Move(ctx context.Context, fileIDs []string, targetFolderID string) (*BulkResult, error) {
	req, err := s.client.newRequest(ctx, "POST", "/files/bulk-move", BulkMoveRequest{
		FileIDs:        fileIDs,
		TargetFolderID: targetFolderID,
	})
	if err != nil {
		return nil, err
	}

	var resp BulkResult
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Download returns a ZIP archive containing the requested files.
func (s *BulkService) Download(ctx context.Context, fileIDs []string) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "POST", "/files/bulk-download", BulkDownloadRequest{FileIDs: fileIDs})
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}
