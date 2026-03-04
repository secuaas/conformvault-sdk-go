// Package conformvault provides a Go SDK for the ConformVault Developer API.
//
// Usage:
//
//	client := conformvault.NewClient("cvk_live_your_api_key")
//	files, err := client.Files.List(ctx, nil)
package conformvault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL.
	DefaultBaseURL = "https://api.conformvault.com/dev/v1"
	// Version is the SDK version.
	Version = "0.5.4"
	// userAgent is the User-Agent header value.
	userAgent = "conformvault-go/" + Version
)

// Client is the ConformVault API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Files                 *FilesService
	Folders               *FoldersService
	ShareLinks            *ShareLinksService
	Signatures            *SignaturesService
	Webhooks              *WebhooksService
	Audit                 *AuditService
	Keys                  *KeysService
	Bulk                  *BulkService
	Versions              *VersionsService
	Search                *SearchService
	Trash                 *TrashService
	ScanReports           *ScanReportsService
	Attestation           *AttestationService
	Transactions          *TransactionsService
	Templates             *TemplatesService
	Batches               *BatchesService
	Metadata              *MetadataService
	Retention             *RetentionService
	LegalHolds            *LegalHoldsService
	Permissions           *PermissionsService
	Comments              *CommentsService
	Quota                 *QuotaService
	RateLimit             *RateLimitService
	UploadSessions        *UploadSessionsService
	Jobs                  *JobsService
	ActivitySubscriptions *ActivitySubscriptionsService
	Policies              *PoliciesService
	Bandwidth             *BandwidthService
	DataExport            *DataExportService
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClient creates a new ConformVault API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Files = &FilesService{client: c}
	c.Folders = &FoldersService{client: c}
	c.ShareLinks = &ShareLinksService{client: c}
	c.Signatures = &SignaturesService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Audit = &AuditService{client: c}
	c.Keys = &KeysService{client: c}
	c.Bulk = &BulkService{client: c}
	c.Versions = &VersionsService{client: c}
	c.Search = &SearchService{client: c}
	c.Trash = &TrashService{client: c}
	c.ScanReports = &ScanReportsService{client: c}
	c.Attestation = &AttestationService{client: c}
	c.Transactions = &TransactionsService{client: c}
	c.Templates = &TemplatesService{client: c}
	c.Batches = &BatchesService{client: c}
	c.Metadata = &MetadataService{client: c}
	c.Retention = &RetentionService{client: c}
	c.LegalHolds = &LegalHoldsService{client: c}
	c.Permissions = &PermissionsService{client: c}
	c.Comments = &CommentsService{client: c}
	c.Quota = &QuotaService{client: c}
	c.RateLimit = &RateLimitService{client: c}
	c.UploadSessions = &UploadSessionsService{client: c}
	c.Jobs = &JobsService{client: c}
	c.ActivitySubscriptions = &ActivitySubscriptionsService{client: c}
	c.Policies = &PoliciesService{client: c}
	c.Bandwidth = &BandwidthService{client: c}
	c.DataExport = &DataExportService{client: c}

	return c
}

// newRequest creates a new HTTP request with authentication headers.
func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes an HTTP request and decodes the JSON response.
func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 60 * time.Second
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if seconds, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}
		return &RateLimitError{
			APIError:   APIError{StatusCode: resp.StatusCode, Message: "rate limited"},
			RetryAfter: retryAfter,
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB max
		if len(body) > 0 {
			if err := json.Unmarshal(body, &apiErr); err != nil {
				apiErr.Message = string(body)
			}
		}
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return &apiErr
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// doRaw executes an HTTP request and returns the raw response body.
func (c *Client) doRaw(req *http.Request) (io.ReadCloser, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		json.NewDecoder(resp.Body).Decode(&apiErr)
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return nil, &apiErr
	}

	return resp.Body, nil
}
