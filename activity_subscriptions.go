package conformvault

import "context"

// ActivitySubscriptionsService manages activity event subscriptions.
type ActivitySubscriptionsService struct {
	client *Client
}

// Subscribe creates a new activity event subscription.
func (s *ActivitySubscriptionsService) Subscribe(ctx context.Context, request CreateActivitySubscriptionRequest) (*ActivitySubscription, error) {
	req, err := s.client.newRequest(ctx, "POST", "/activity-subscriptions", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[ActivitySubscription]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all activity event subscriptions.
func (s *ActivitySubscriptionsService) List(ctx context.Context) ([]ActivitySubscription, error) {
	req, err := s.client.newRequest(ctx, "GET", "/activity-subscriptions", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[ActivitySubscription]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Unsubscribe removes an activity event subscription.
func (s *ActivitySubscriptionsService) Unsubscribe(ctx context.Context, subscriptionID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/activity/subscriptions/"+subscriptionID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
