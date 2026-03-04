package conformvault

import "context"

// BandwidthService provides bandwidth analytics.
type BandwidthService struct {
	client *Client
}

// BandwidthSummary contains bandwidth usage summary.
type BandwidthSummary struct {
	TotalUploadBytes   int64  `json:"total_upload_bytes"`
	TotalDownloadBytes int64  `json:"total_download_bytes"`
	Period             string `json:"period"`
}

// DailyBandwidthStats contains daily bandwidth statistics.
type DailyBandwidthStats struct {
	Date          string `json:"date"`
	UploadBytes   int64  `json:"upload_bytes"`
	DownloadBytes int64  `json:"download_bytes"`
}

// GetSummary returns the bandwidth usage summary.
func (s *BandwidthService) GetSummary(ctx context.Context) (*BandwidthSummary, error) {
	req, err := s.client.newRequest(ctx, "GET", "/bandwidth", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[BandwidthSummary]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetDaily returns daily bandwidth statistics.
func (s *BandwidthService) GetDaily(ctx context.Context) ([]DailyBandwidthStats, error) {
	req, err := s.client.newRequest(ctx, "GET", "/bandwidth/daily", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[DailyBandwidthStats]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
