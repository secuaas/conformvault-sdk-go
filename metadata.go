package conformvault

import "context"

// MetadataService handles file metadata and tags.
type MetadataService struct {
	client *Client
}

// AddTags adds tags to a file.
func (s *MetadataService) AddTags(ctx context.Context, fileID string, request AddTagsRequest) ([]FileTag, error) {
	req, err := s.client.newRequest(ctx, "POST", "/files/"+fileID+"/tags", request)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[FileTag]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// RemoveTag removes a tag from a file.
func (s *MetadataService) RemoveTag(ctx context.Context, fileID, tag string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/files/"+fileID+"/tags/"+tag, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetTags returns all tags for a file.
func (s *MetadataService) GetTags(ctx context.Context, fileID string) ([]FileTag, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/tags", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[FileTag]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListByTag returns all files with a specific tag.
func (s *MetadataService) ListByTag(ctx context.Context, tag string) ([]File, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/by-tag/"+tag, nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[File]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// SetMetadata sets metadata key-value pairs for a file.
func (s *MetadataService) SetMetadata(ctx context.Context, fileID string, request SetMetadataRequest) (map[string]string, error) {
	req, err := s.client.newRequest(ctx, "PATCH", "/files/"+fileID+"/metadata", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[map[string]string]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetMetadata returns all metadata for a file.
func (s *MetadataService) GetMetadata(ctx context.Context, fileID string) (map[string]string, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/metadata", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[map[string]string]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeleteMetadataKey removes a metadata key from a file.
func (s *MetadataService) DeleteMetadataKey(ctx context.Context, fileID, key string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/files/"+fileID+"/metadata/"+key, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
