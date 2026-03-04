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

// ListDeliveries returns delivery attempts for a webhook.
func (s *WebhooksService) ListDeliveries(ctx context.Context, webhookID string) ([]WebhookDelivery, error) {
	req, err := s.client.newRequest(ctx, "GET", "/webhooks/"+webhookID+"/deliveries", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[WebhookDelivery]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetDelivery returns a specific delivery attempt.
func (s *WebhooksService) GetDelivery(ctx context.Context, webhookID, deliveryID string) (*WebhookDelivery, error) {
	req, err := s.client.newRequest(ctx, "GET", "/webhooks/"+webhookID+"/deliveries/"+deliveryID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[WebhookDelivery]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ReplayDelivery replays a failed delivery attempt.
func (s *WebhooksService) ReplayDelivery(ctx context.Context, webhookID, deliveryID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/webhooks/"+webhookID+"/deliveries/"+deliveryID+"/replay", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// ReEnable re-enables a disabled webhook endpoint.
func (s *WebhooksService) ReEnable(ctx context.Context, webhookID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/webhooks/"+webhookID+"/enable", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
