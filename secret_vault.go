package conformvault

import "context"

// SecretVaultService manages ephemeral secrets.
type SecretVaultService struct {
	client *Client
}

// Create creates a new secret.
func (s *SecretVaultService) Create(ctx context.Context, request CreateSecretRequest) (*Secret, error) {
	req, err := s.client.newRequest(ctx, "POST", "/secrets", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Secret]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all secrets.
func (s *SecretVaultService) List(ctx context.Context, opts *SecretListOptions) ([]Secret, error) {
	path := "/secrets"
	if opts != nil {
		sep := "?"
		if opts.Limit > 0 {
			path += sep + "limit=" + itoa(opts.Limit)
			sep = "&"
		}
		if opts.Offset > 0 {
			path += sep + "offset=" + itoa(opts.Offset)
		}
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Secret]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single secret by ID.
func (s *SecretVaultService) Get(ctx context.Context, id string) (*Secret, error) {
	req, err := s.client.newRequest(ctx, "GET", "/secrets/"+id, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Secret]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a secret by ID.
func (s *SecretVaultService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/secrets/"+id, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
