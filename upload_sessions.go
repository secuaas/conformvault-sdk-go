package conformvault

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// UploadSessionsService manages chunked upload sessions.
type UploadSessionsService struct {
	client *Client
}

// Create creates a new chunked upload session.
func (s *UploadSessionsService) Create(ctx context.Context, request CreateUploadSessionRequest) (*UploadSession, error) {
	req, err := s.client.newRequest(ctx, "POST", "/upload-sessions", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[UploadSession]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UploadChunk uploads a single chunk of data to an upload session.
func (s *UploadSessionsService) UploadChunk(ctx context.Context, sessionID string, chunkNumber int, data io.Reader) error {
	path := fmt.Sprintf("/upload-sessions/%s/chunks/%d", sessionID, chunkNumber)
	url := s.client.baseURL + path
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, data)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("Content-Type", "application/octet-stream")
	return s.client.do(httpReq, nil)
}

// GetStatus returns the current status of an upload session.
func (s *UploadSessionsService) GetStatus(ctx context.Context, sessionID string) (*UploadSession, error) {
	req, err := s.client.newRequest(ctx, "GET", "/upload-sessions/"+sessionID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[UploadSession]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Complete finalizes an upload session and returns the created file.
func (s *UploadSessionsService) Complete(ctx context.Context, sessionID string) (*File, error) {
	req, err := s.client.newRequest(ctx, "POST", "/upload-sessions/"+sessionID+"/complete", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[File]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Cancel cancels an upload session.
func (s *UploadSessionsService) Cancel(ctx context.Context, sessionID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/upload-sessions/"+sessionID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
