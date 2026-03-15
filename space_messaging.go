package conformvault

import (
	"context"
	"fmt"
)

// SpaceMessagingService manages messages within spaces.
type SpaceMessagingService struct {
	client *Client
}

// spaceMessageListResponse wraps paginated space messages with a total count.
type spaceMessageListResponse struct {
	Data  []SpaceMessage `json:"data"`
	Total int            `json:"total"`
}

// SendMessage sends a message in a space.
func (s *SpaceMessagingService) SendMessage(ctx context.Context, spaceID string, request CreateSpaceMessageRequest) (*SpaceMessage, error) {
	req, err := s.client.newRequest(ctx, "POST", "/spaces/"+spaceID+"/messages", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[SpaceMessage]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListMessages returns messages in a space with optional pagination.
func (s *SpaceMessagingService) ListMessages(ctx context.Context, spaceID string, opts *MessageListOptions) ([]SpaceMessage, int, error) {
	path := "/spaces/" + spaceID + "/messages"
	if opts != nil {
		sep := "?"
		if opts.Limit > 0 {
			path += sep + fmt.Sprintf("limit=%d", opts.Limit)
			sep = "&"
		}
		if opts.Offset > 0 {
			path += sep + fmt.Sprintf("offset=%d", opts.Offset)
		}
	}
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, 0, err
	}
	var resp spaceMessageListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Data, resp.Total, nil
}

// GetReplies returns replies to a specific message.
func (s *SpaceMessagingService) GetReplies(ctx context.Context, spaceID, messageID string) ([]SpaceMessage, error) {
	req, err := s.client.newRequest(ctx, "GET", "/spaces/"+spaceID+"/messages/"+messageID+"/replies", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[SpaceMessage]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// MarkRead marks a message as read.
func (s *SpaceMessagingService) MarkRead(ctx context.Context, spaceID, messageID string) error {
	req, err := s.client.newRequest(ctx, "POST", "/spaces/"+spaceID+"/messages/"+messageID+"/read", nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// DeleteMessage deletes a message.
func (s *SpaceMessagingService) DeleteMessage(ctx context.Context, spaceID, messageID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/spaces/"+spaceID+"/messages/"+messageID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
