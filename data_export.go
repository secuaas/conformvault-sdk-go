package conformvault

import "context"

// DataExportService handles GDPR/Loi 25 user data exports.
type DataExportService struct {
	client *Client
}

// UserDataExport represents an exported user data package.
type UserDataExport struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	URL       string `json:"url,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// Export requests a data export for a user (GDPR/Loi 25 compliance).
func (s *DataExportService) Export(ctx context.Context, userID string) (*UserDataExport, error) {
	req, err := s.client.newRequest(ctx, "GET", "/users/"+userID+"/export", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[UserDataExport]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
