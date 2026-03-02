package conformvault

import (
	"context"
	"fmt"
	"time"
)

// VersionsService handles file version operations.
type VersionsService struct {
	client *Client
}

// FileVersion represents a version of a file.
type FileVersion struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	Version   int       `json:"version"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// List returns all versions of a file.
func (s *VersionsService) List(ctx context.Context, fileID string) ([]FileVersion, error) {
	req, err := s.client.newRequest(ctx, "GET", fmt.Sprintf("/files/%s/versions", fileID), nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[FileVersion]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a specific version of a file.
func (s *VersionsService) Get(ctx context.Context, fileID, versionID string) (*FileVersion, error) {
	req, err := s.client.newRequest(ctx, "GET", fmt.Sprintf("/files/%s/versions/%s", fileID, versionID), nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[FileVersion]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Restore restores an old version as the current version.
func (s *VersionsService) Restore(ctx context.Context, fileID, versionID string) error {
	req, err := s.client.newRequest(ctx, "POST", fmt.Sprintf("/files/%s/versions/%s/restore", fileID, versionID), nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Delete permanently deletes a file version.
func (s *VersionsService) Delete(ctx context.Context, fileID, versionID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", fmt.Sprintf("/files/%s/versions/%s", fileID, versionID), nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
