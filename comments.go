package conformvault

import "context"

// CommentsService manages file comments.
type CommentsService struct {
	client *Client
}

// Create creates a new comment on a file.
func (s *CommentsService) Create(ctx context.Context, request CreateCommentRequest) (*Comment, error) {
	req, err := s.client.newRequest(ctx, "POST", "/comments", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns all comments for a file.
func (s *CommentsService) List(ctx context.Context, fileID string) ([]Comment, error) {
	req, err := s.client.newRequest(ctx, "GET", "/comments?file_id="+fileID, nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single comment by ID.
func (s *CommentsService) Get(ctx context.Context, commentID string) (*Comment, error) {
	req, err := s.client.newRequest(ctx, "GET", "/comments/"+commentID, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing comment.
func (s *CommentsService) Update(ctx context.Context, commentID string, request UpdateCommentRequest) (*Comment, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/comments/"+commentID, request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a comment by ID.
func (s *CommentsService) Delete(ctx context.Context, commentID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/comments/"+commentID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetReplies returns all replies to a comment.
func (s *CommentsService) GetReplies(ctx context.Context, commentID string) ([]Comment, error) {
	req, err := s.client.newRequest(ctx, "GET", "/comments/"+commentID+"/replies", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
