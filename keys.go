package conformvault

import (
	"context"
	"fmt"
)

// KeysService handles API key self-management.
type KeysService struct {
	client *Client
}

// List returns all API keys for the organization.
func (s *KeysService) List(ctx context.Context, opts *KeyListOptions) ([]APIKey, error) {
	path := "/keys"
	if opts != nil {
		sep := "?"
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

	var resp ListResponse[APIKey]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Create creates a new API key. The full key is returned once.
func (s *KeysService) Create(ctx context.Context, r CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/keys", r)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[CreateAPIKeyResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get returns a single API key by ID.
func (s *KeysService) Get(ctx context.Context, keyID string) (*APIKey, error) {
	req, err := s.client.newRequest(ctx, "GET", "/keys/"+keyID, nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[APIKey]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Revoke revokes (deletes) an API key.
func (s *KeysService) Revoke(ctx context.Context, keyID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/keys/"+keyID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Rotate rotates an API key, returning a new key value.
func (s *KeysService) Rotate(ctx context.Context, keyID string) (*CreateAPIKeyResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/keys/"+keyID+"/rotate", nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[CreateAPIKeyResponse]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// KeyRevocationStatus represents the revocation status of an API key.
type KeyRevocationStatus struct {
	KeyID     string `json:"key_id"`
	Revoked   bool   `json:"revoked"`
	RevokedAt string `json:"revoked_at,omitempty"`
}

// InstantRevoke instantly revokes an API key via Redis.
func (s *KeysService) InstantRevoke(ctx context.Context, keyID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/api-keys/"+keyID+"/revoke", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetRevocationStatus checks the revocation status of an API key.
func (s *KeysService) GetRevocationStatus(ctx context.Context, keyID string) (*KeyRevocationStatus, error) {
	req, err := s.client.newRequest(ctx, "GET", "/api-keys/"+keyID+"/revocation-status", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[KeyRevocationStatus]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
