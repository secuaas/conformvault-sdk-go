package conformvault

import "time"

// --- Files ---

// File represents a stored file.
type File struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	FolderID     *string   `json:"folder_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// FileListOptions are query parameters for listing files.
type FileListOptions struct {
	FolderID *string `json:"folder_id,omitempty"`
	Page     int     `json:"page,omitempty"`
	Limit    int     `json:"limit,omitempty"`
}

// UploadResult is the response from uploading a file.
type UploadResult struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

// --- Folders ---

// Folder represents a folder in the file tree.
type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	ParentID  *string   `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// FolderListOptions are query parameters for listing folders.
type FolderListOptions struct {
	ParentID *string `json:"parent_id,omitempty"`
	Page     int     `json:"page,omitempty"`
	Limit    int     `json:"limit,omitempty"`
}

// CreateFolderRequest is the input for creating a folder.
type CreateFolderRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
}

// --- Share Links ---

// ShareLink represents a share link.
type ShareLink struct {
	ID          string    `json:"id"`
	Token       string    `json:"token"`
	URL         string    `json:"url"`
	Type        string    `json:"type"` // "download" or "upload"
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// ShareLinkListOptions are query parameters for listing share links.
type ShareLinkListOptions struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// CreateShareLinkRequest is the input for creating a share link.
type CreateShareLinkRequest struct {
	FileID    *string `json:"file_id,omitempty"`
	FolderID  *string `json:"folder_id,omitempty"`
	Type      string  `json:"type"`
	ExpiresIn int     `json:"expires_in,omitempty"` // seconds
	Password  *string `json:"password,omitempty"`
}

// --- Signatures ---

// SignatureEnvelope represents a signature envelope.
type SignatureEnvelope struct {
	ID           string    `json:"id"`
	Provider     string    `json:"provider"`
	Status       string    `json:"status"`
	Subject      string    `json:"subject"`
	Message      *string   `json:"message,omitempty"`
	SourceFileID *string   `json:"source_file_id,omitempty"`
	SignedFileID *string   `json:"signed_file_id,omitempty"`
	ExpiryDays   int       `json:"expiry_days"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateSignatureRequest is the input for creating a signature envelope.
type CreateSignatureRequest struct {
	FileID     string                 `json:"file_id"`
	Subject    string                 `json:"subject"`
	Message    *string                `json:"message,omitempty"`
	Signers    []CreateSignatureSigner `json:"signers"`
	ExpiryDays int                    `json:"expiry_days,omitempty"`
}

// CreateSignatureSigner is a signer in a create signature request.
type CreateSignatureSigner struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role,omitempty"`
	SignOrder int    `json:"sign_order,omitempty"`
}

// --- Webhooks ---

// WebhookEndpoint represents a registered webhook endpoint.
type WebhookEndpoint struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Events      []string  `json:"events"`
	Environment string    `json:"environment"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterWebhookRequest is the input for registering a webhook.
type RegisterWebhookRequest struct {
	URL         string   `json:"url"`
	Events      []string `json:"events,omitempty"`
	Environment string   `json:"environment,omitempty"`
}

// RegisterWebhookResponse includes the signing secret (shown once).
type RegisterWebhookResponse struct {
	WebhookEndpoint
	Secret string `json:"secret"`
}

// --- Audit ---

// AuditEntry represents an audit log entry.
type AuditEntry struct {
	ID           string    `json:"id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Details      any       `json:"details,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditListOptions are query parameters for listing audit entries.
type AuditListOptions struct {
	EventType string `json:"event_type,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// --- API Keys ---

// APIKey represents a developer API key.
type APIKey struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Environment string    `json:"environment"`
	Scopes      []string  `json:"scopes"`
	ExpiresAt   *string   `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAPIKeyRequest is the input for creating an API key.
type CreateAPIKeyRequest struct {
	Name        string   `json:"name"`
	Environment string   `json:"environment"`
	Scopes      []string `json:"scopes"`
}

// CreateAPIKeyResponse includes the full key (shown once).
type CreateAPIKeyResponse struct {
	APIKey
	Key string `json:"key"`
}

// --- Generic ---

// ListResponse wraps a paginated response.
type ListResponse[T any] struct {
	Data []T `json:"data"`
}

// DataResponse wraps a single-item response.
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// MessageResponse wraps a message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// --- Scan Reports ---

// FileScanReport represents the antivirus scan result for a file.
type FileScanReport struct {
	ID             string  `json:"id"`
	FileID         string  `json:"file_id"`
	OrganizationID string  `json:"organization_id"`
	ScanEngine     string  `json:"scan_engine"`
	EngineVersion  *string `json:"engine_version,omitempty"`
	ScanStatus     string  `json:"scan_status"`
	ThreatName     *string `json:"threat_name,omitempty"`
	FileSize       *int64  `json:"file_size,omitempty"`
	MimeType       *string `json:"mime_type,omitempty"`
	ScanDurationMs *int    `json:"scan_duration_ms,omitempty"`
	ScannedAt      string  `json:"scanned_at"`
}

// ScanReportListOptions are query parameters for listing scan reports.
type ScanReportListOptions struct {
	Limit  int
	Offset int
}

// ScanReportListResponse is the paginated response for listing scan reports.
type ScanReportListResponse struct {
	Reports []FileScanReport `json:"reports"`
	Total   int              `json:"total"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
}

// FileScanSummary contains aggregate scan statistics for an organization.
type FileScanSummary struct {
	TotalScans    int    `json:"total_scans"`
	CleanCount    int    `json:"clean_count"`
	InfectedCount int    `json:"infected_count"`
	ErrorCount    int    `json:"error_count"`
	SkippedCount  int    `json:"skipped_count"`
	ScanEngine    string `json:"scan_engine"`
	EngineVersion string `json:"engine_version"`
}
