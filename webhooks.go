package conformvault

import "context"

// WebhooksService handles webhook endpoint management.
type WebhooksService struct {
	client *Client
}

// List returns all registered webhook endpoints.
func (s *WebhooksService) List(ctx context.Context) ([]WebhookEndpoint, error) {
	req, err := s.client.newRequest(ctx, "GET", "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[WebhookEndpoint]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Register creates a new webhook endpoint. The signing secret is returned once.
func (s *WebhooksService) Register(ctx context.Context, r RegisterWebhookRequest) (*RegisterWebhookResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/webhooks", r)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data    RegisterWebhookResponse `json:"data"`
		Message string                  `json:"message"`
	}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a webhook endpoint.
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/webhooks/"+webhookID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Test sends a test event to a webhook endpoint.
func (s *WebhooksService) Test(ctx context.Context, webhookID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/webhooks/"+webhookID+"/test", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
