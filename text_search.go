package conformvault

import (
	"context"
	"fmt"
	"net/url"
)

// TextSearchService handles full-text content search (OCR + extracted text).
type TextSearchService struct {
	client *Client
}

// Search performs a full-text search across indexed document content.
func (s *TextSearchService) Search(ctx context.Context, opts TextSearchOptions) (*TextSearchResponse, error) {
	path := fmt.Sprintf("/text-search?q=%s", url.QueryEscape(opts.Query))

	if opts.SpaceID != nil {
		path += "&space_id=" + url.QueryEscape(*opts.SpaceID)
	}
	if opts.Language != "" {
		path += "&lang=" + url.QueryEscape(opts.Language)
	}
	if opts.Limit > 0 {
		path += fmt.Sprintf("&limit=%d", opts.Limit)
	}
	if opts.Offset > 0 {
		path += fmt.Sprintf("&offset=%d", opts.Offset)
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp TextSearchResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetIndexStats returns indexing statistics for the organization.
func (s *TextSearchService) GetIndexStats(ctx context.Context) (map[string]any, error) {
	req, err := s.client.newRequest(ctx, "GET", "/text-search/stats", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data map[string]any `json:"data"`
	}
	if err := s.client.do(req, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}

// GetFileIndexStatus returns the index status for a specific file.
func (s *TextSearchService) GetFileIndexStatus(ctx context.Context, fileID string) (*DocumentTextIndex, error) {
	path := fmt.Sprintf("/text-search/files/%s/status", url.PathEscape(fileID))

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data DocumentTextIndex `json:"data"`
	}
	if err := s.client.do(req, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// ReindexFile re-queues a file for text indexing.
func (s *TextSearchService) ReindexFile(ctx context.Context, fileID string) error {
	path := fmt.Sprintf("/text-search/files/%s/reindex", url.PathEscape(fileID))

	req, err := s.client.newRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

// RequeueSkippedOCR re-queues all previously skipped OCR entries for processing.
func (s *TextSearchService) RequeueSkippedOCR(ctx context.Context) (*RequeueOCRResponse, error) {
	req, err := s.client.newRequest(ctx, "POST", "/text-search/requeue-ocr", nil)
	if err != nil {
		return nil, err
	}

	var resp RequeueOCRResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
