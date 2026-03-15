package conformvault

import (
	"context"
	"io"
	"net/url"
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

// List returns all signature envelopes for the organization.
func (s *SignaturesService) List(ctx context.Context) ([]SignatureEnvelope, error) {
	req, err := s.client.newRequest(ctx, "GET", "/signatures", nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[[]SignatureEnvelope]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DownloadAuditTrail downloads the audit trail (completion certificate) PDF.
func (s *SignaturesService) DownloadAuditTrail(ctx context.Context, envelopeID string) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "GET", "/signatures/"+envelopeID+"/audit-trail", nil)
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

// AnalyzePDF analyzes a PDF for signature field placement.
func (s *SignaturesService) AnalyzePDF(ctx context.Context, req AnalyzePDFRequest) (*PDFAnalysisResult, error) {
	r, err := s.client.newRequest(ctx, "POST", "/signatures/analyze", req)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[PDFAnalysisResult]
	if err := s.client.do(r, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// PreviewPDF streams a decrypted PDF for signature placement preview.
func (s *SignaturesService) PreviewPDF(ctx context.Context, fileID string) (io.ReadCloser, error) {
	r, err := s.client.newRequest(ctx, "GET", "/signatures/preview-pdf?file_id="+url.QueryEscape(fileID), nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(r)
}

// Delegate delegates a signer's signature to another person.
func (s *SignaturesService) Delegate(ctx context.Context, envelopeID, signerID string, request DelegateSignRequest) error {
	req, err := s.client.newRequest(ctx, "POST", "/signatures/"+envelopeID+"/signers/"+signerID+"/delegate", request)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetEmbeddedSignLink retrieves an embedded signing link for a signer.
func (s *SignaturesService) GetEmbeddedSignLink(ctx context.Context, envelopeID, signerEmail string, redirectURL ...string) (*EmbeddedSignLinkResponse, error) {
	path := "/signatures/" + envelopeID + "/embed-sign?signer_email=" + url.QueryEscape(signerEmail)
	if len(redirectURL) > 0 && redirectURL[0] != "" {
		path += "&redirect_url=" + url.QueryEscape(redirectURL[0])
	}
	r, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp EmbeddedSignLinkResponse
	if err := s.client.do(r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
