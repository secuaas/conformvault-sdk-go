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

// --- Webhook Delivery types ---

// WebhookDelivery represents a webhook delivery attempt.
type WebhookDelivery struct {
	ID           string     `json:"id"`
	WebhookID    string     `json:"webhook_id"`
	EventType    string     `json:"event_type"`
	Status       string     `json:"status"`
	HTTPStatus   int        `json:"http_status"`
	RequestBody  string     `json:"request_body,omitempty"`
	ResponseBody string     `json:"response_body,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
}

// --- Audit extended types ---

// AuditSearchOptions are query parameters for searching audit entries.
type AuditSearchOptions struct {
	Query     string `json:"query,omitempty"`
	EventType string `json:"event_type,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// AuditExportOptions are query parameters for exporting audit logs.
type AuditExportOptions struct {
	Format    string `json:"format,omitempty"` // json, csv
	EventType string `json:"event_type,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
}

// AuditStats contains aggregated audit statistics.
type AuditStats struct {
	TotalEvents  int            `json:"total_events"`
	EventsByType map[string]int `json:"events_by_type"`
	EventsByDay  map[string]int `json:"events_by_day"`
}

// AuditAnomaly represents a detected audit anomaly.
type AuditAnomaly struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	DetectedAt  time.Time `json:"detected_at"`
}

// --- File Metadata & Tags types ---

// FileTag represents a tag attached to a file.
type FileTag struct {
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

// AddTagsRequest is the input for adding tags to a file.
type AddTagsRequest struct {
	Tags []string `json:"tags"`
}

// FileMetadataEntry represents a single metadata key-value pair.
type FileMetadataEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetMetadataRequest is the input for setting metadata on a file.
type SetMetadataRequest struct {
	Metadata map[string]string `json:"metadata"`
}

// --- Retention Policy types ---

// RetentionPolicy represents a data retention policy.
type RetentionPolicy struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	RetentionDays int       `json:"retention_days"`
	AutoDelete    bool      `json:"auto_delete"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateRetentionPolicyRequest is the input for creating a retention policy.
type CreateRetentionPolicyRequest struct {
	Name          string `json:"name"`
	RetentionDays int    `json:"retention_days"`
	AutoDelete    bool   `json:"auto_delete"`
}

// UpdateRetentionPolicyRequest is the input for updating a retention policy.
type UpdateRetentionPolicyRequest struct {
	Name          *string `json:"name,omitempty"`
	RetentionDays *int    `json:"retention_days,omitempty"`
	AutoDelete    *bool   `json:"auto_delete,omitempty"`
}

// --- Legal Hold types ---

// LegalHold represents a legal hold on files.
type LegalHold struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	FileCount   int        `json:"file_count"`
	CreatedAt   time.Time  `json:"created_at"`
	ReleasedAt  *time.Time `json:"released_at,omitempty"`
}

// CreateLegalHoldRequest is the input for creating a legal hold.
type CreateLegalHoldRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// AddLegalHoldFilesRequest is the input for adding files to a legal hold.
type AddLegalHoldFilesRequest struct {
	FileIDs []string `json:"file_ids"`
}

// LegalHoldFile represents a file under legal hold.
type LegalHoldFile struct {
	FileID  string    `json:"file_id"`
	AddedAt time.Time `json:"added_at"`
}

// --- Folder Permission types ---

// FolderPermission represents a permission granted on a folder.
type FolderPermission struct {
	FolderID   string    `json:"folder_id"`
	UserID     string    `json:"user_id"`
	Permission string    `json:"permission"`
	GrantedAt  time.Time `json:"granted_at"`
}

// SetFolderPermissionRequest is the input for setting a folder permission.
type SetFolderPermissionRequest struct {
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}

// --- Comment types ---

// Comment represents a comment on a file.
type Comment struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	ParentID  *string   `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCommentRequest is the input for creating a comment.
type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id,omitempty"`
}

// UpdateCommentRequest is the input for updating a comment.
type UpdateCommentRequest struct {
	Content string `json:"content"`
}

// --- Quota types ---

// QuotaInfo contains storage quota information.
type QuotaInfo struct {
	UsedBytes    int64 `json:"used_bytes"`
	TotalBytes   int64 `json:"total_bytes"`
	FileCount    int   `json:"file_count"`
	MaxFileCount int   `json:"max_file_count"`
}

// --- Rate Limit types ---

// RateLimitInfo contains rate limit status information.
type RateLimitInfo struct {
	RequestsPerMinute int    `json:"requests_per_minute"`
	RequestsRemaining int    `json:"requests_remaining"`
	ResetAt           string `json:"reset_at"`
}

// --- Upload Session types ---

// UploadSession represents a chunked upload session.
type UploadSession struct {
	ID             string    `json:"id"`
	Filename       string    `json:"filename"`
	TotalSize      int64     `json:"total_size"`
	ChunkSize      int64     `json:"chunk_size"`
	ChunksUploaded int       `json:"chunks_uploaded"`
	TotalChunks    int       `json:"total_chunks"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// CreateUploadSessionRequest is the input for creating a chunked upload session.
type CreateUploadSessionRequest struct {
	Filename    string  `json:"filename"`
	TotalSize   int64   `json:"total_size"`
	ContentType string  `json:"content_type,omitempty"`
	FolderID    *string `json:"folder_id,omitempty"`
}

// --- Job types ---

// Job represents a background job.
type Job struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	Result      any        `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// CreateJobRequest is the input for creating a background job.
type CreateJobRequest struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params,omitempty"`
}

// --- Activity Subscription types ---

// ActivitySubscription represents an activity event subscription.
type ActivitySubscription struct {
	ID          string    `json:"id"`
	EventTypes  []string  `json:"event_types"`
	CallbackURL string    `json:"callback_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateActivitySubscriptionRequest is the input for creating an activity subscription.
type CreateActivitySubscriptionRequest struct {
	EventTypes  []string `json:"event_types"`
	CallbackURL string   `json:"callback_url"`
}

// --- IP Policy types ---

// IPPolicy represents an IP restriction policy.
type IPPolicy struct {
	Enabled    bool     `json:"enabled"`
	AllowedIPs []string `json:"allowed_ips"`
	DeniedIPs  []string `json:"denied_ips"`
}

// SetIPPolicyRequest is the input for updating the IP restriction policy.
type SetIPPolicyRequest struct {
	Enabled    bool     `json:"enabled"`
	AllowedIPs []string `json:"allowed_ips,omitempty"`
	DeniedIPs  []string `json:"denied_ips,omitempty"`
}

// --- MFA Policy types ---

// MFAPolicy represents a multi-factor authentication policy.
type MFAPolicy struct {
	Enabled     bool     `json:"enabled"`
	RequiredFor []string `json:"required_for"`
}

// SetMFAPolicyRequest is the input for updating the MFA policy.
type SetMFAPolicyRequest struct {
	Enabled     bool     `json:"enabled"`
	RequiredFor []string `json:"required_for,omitempty"`
}

// --- Encryption Salt types ---

// EncryptionSalt represents the encryption salt configuration.
type EncryptionSalt struct {
	Salt string `json:"salt"`
}

// SetEncryptionSaltRequest is the input for updating the encryption salt.
type SetEncryptionSaltRequest struct {
	Salt string `json:"salt"`
}
