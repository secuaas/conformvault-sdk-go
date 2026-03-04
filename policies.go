package conformvault

import "context"

// PoliciesService manages security policies and encryption configuration.
type PoliciesService struct {
	client *Client
}

// GetIPPolicy returns the current IP restriction policy.
func (s *PoliciesService) GetIPPolicy(ctx context.Context) (*IPPolicy, error) {
	req, err := s.client.newRequest(ctx, "GET", "/policies/ip", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[IPPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SetIPPolicy updates the IP restriction policy.
func (s *PoliciesService) SetIPPolicy(ctx context.Context, request SetIPPolicyRequest) (*IPPolicy, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/policies/ip", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[IPPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetMFAPolicy returns the current MFA policy.
func (s *PoliciesService) GetMFAPolicy(ctx context.Context) (*MFAPolicy, error) {
	req, err := s.client.newRequest(ctx, "GET", "/policies/mfa", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[MFAPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SetMFAPolicy updates the MFA policy.
func (s *PoliciesService) SetMFAPolicy(ctx context.Context, request SetMFAPolicyRequest) (*MFAPolicy, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/policies/mfa", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[MFAPolicy]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetEncryptionSalt returns the current encryption salt.
func (s *PoliciesService) GetEncryptionSalt(ctx context.Context) (*EncryptionSalt, error) {
	req, err := s.client.newRequest(ctx, "GET", "/policies/encryption-salt", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[EncryptionSalt]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SetEncryptionSalt updates the encryption salt.
func (s *PoliciesService) SetEncryptionSalt(ctx context.Context, request SetEncryptionSaltRequest) (*EncryptionSalt, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/policies/encryption-salt", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[EncryptionSalt]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
