package conformvault

import (
	"context"
	"io"
)

// AttestationService handles compliance attestation operations.
type AttestationService struct {
	client *Client
}

// GenerateLoi25 returns a reader for the Loi 25 compliance attestation PDF.
func (s *AttestationService) GenerateLoi25(ctx context.Context) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "GET", "/attestation/loi25", nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}
