package conformvault

import "context"

// PermissionsService manages folder permissions.
type PermissionsService struct {
	client *Client
}

// Set sets a permission on a folder for a user.
func (s *PermissionsService) Set(ctx context.Context, folderID string, request SetFolderPermissionRequest) (*FolderPermission, error) {
	req, err := s.client.newRequest(ctx, "POST", "/folders/"+folderID+"/permissions", request)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[FolderPermission]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get returns all permissions for a folder.
func (s *PermissionsService) Get(ctx context.Context, folderID string) ([]FolderPermission, error) {
	req, err := s.client.newRequest(ctx, "GET", "/folders/"+folderID+"/permissions", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[FolderPermission]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Revoke revokes a user's permission on a folder.
func (s *PermissionsService) Revoke(ctx context.Context, folderID, userID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/folders/"+folderID+"/permissions/"+userID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}
