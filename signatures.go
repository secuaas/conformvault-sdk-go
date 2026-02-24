package conformvault

import (
	"context"
	"io"
)

// SignaturesService handles electronic signature operations.
type SignaturesService struct {
	client *Client
}

// Create creates a new signature envelope.
func (s *SignaturesService) Create(ctx context.Context, r CreateSignatureRequest) (*SignatureEnvelope, error) {
	req, err := s.client.newRequest(ctx, "POST", "/signatures", r)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[SignatureEnvelope]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetStatus retrieves the current status of a signature envelope.
func (s *SignaturesService) GetStatus(ctx context.Context, envelopeID string) (*SignatureEnvelope, error) {
	req, err := s.client.newRequest(ctx, "GET", "/signatures/"+envelopeID, nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[SignatureEnvelope]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DownloadSigned downloads the completed signed document.
func (s *SignaturesService) DownloadSigned(ctx context.Context, envelopeID string) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "GET", "/signatures/"+envelopeID+"/download", nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}

// Revoke cancels a pending signature envelope.
func (s *SignaturesService) Revoke(ctx context.Context, envelopeID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/signatures/"+envelopeID+"/revoke", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
