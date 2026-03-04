# ConformVault Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/secuaas/conformvault-sdk-go.svg)](https://pkg.go.dev/github.com/secuaas/conformvault-sdk-go)
[![CI](https://github.com/secuaas/conformvault-sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/secuaas/conformvault-sdk-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Official Go SDK for the [ConformVault](https://conformvault.com) Developer API -- secure file storage, electronic signatures, and compliance automation.

## Installation

```bash
go get github.com/secuaas/conformvault-sdk-go
```

**Requirements:** Go 1.21+

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cv "github.com/secuaas/conformvault-sdk-go"
)

func main() {
	client := cv.NewClient(os.Getenv("CONFORMVAULT_API_KEY"))
	ctx := context.Background()

	// List files
	files, err := client.Files.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Printf("%s -- %d bytes\n", f.OriginalName, f.Size)
	}

	// Upload a file
	file, _ := os.Open("document.pdf")
	defer file.Close()
	result, err := client.Files.Upload(ctx, file, "document.pdf", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded: %s\n", result.ID)

	// Download a file
	reader, err := client.Files.Download(ctx, result.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	// ... write reader to file or process stream
}
```

## Configuration

The client accepts functional options to customize its behavior:

```go
import (
	"net/http"
	"time"

	cv "github.com/secuaas/conformvault-sdk-go"
)

client := cv.NewClient("cvk_live_your_api_key",
	// Custom API base URL
	cv.WithBaseURL("https://custom-api.example.com/dev/v1"),

	// Custom HTTP client (timeouts, proxies, transport, etc.)
	cv.WithHTTPClient(&http.Client{
		Timeout: 60 * time.Second,
	}),
)
```

### Available Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithBaseURL(url)` | Override the API base URL | `https://api.conformvault.com/dev/v1` |
| `WithHTTPClient(client)` | Provide a custom `*http.Client` | 30-second timeout client |

## API Base URL

| Environment | URL |
|-------------|-----|
| Production | `https://api.conformvault.com/dev/v1` (default) |

## Authentication

All requests use Bearer token authentication:

```
Authorization: Bearer cvk_live_xxx
```

API keys prefixed with `cvk_live_` are for production; `cvk_test_` for sandbox.

## Services

The SDK provides 27 service clients matching the ConformVault Developer API:

| Service | Methods | Accessor |
|---------|---------|----------|
| **Files** | `List`, `Get`, `Upload`, `Download`, `Delete`, `GetThumbnail`, `GetScanReport` | `client.Files` |
| **Folders** | `List`, `Get`, `Create`, `Delete` | `client.Folders` |
| **ShareLinks** | `List`, `Create`, `Delete` | `client.ShareLinks` |
| **Signatures** | `Create`, `GetStatus`, `DownloadSigned`, `Revoke` | `client.Signatures` |
| **Webhooks** | `List`, `Register`, `Delete`, `Test` | `client.Webhooks` |
| **Audit** | `List` | `client.Audit` |
| **Keys** | `List`, `Create`, `Get`, `Revoke`, `Rotate` | `client.Keys` |
| **Bulk** | `Delete`, `Move`, `Download` | `client.Bulk` |
| **Versions** | `List`, `Get`, `Restore`, `Delete` | `client.Versions` |
| **Search** | `Search` | `client.Search` |
| **Trash** | `List`, `Restore`, `Delete`, `Empty` | `client.Trash` |
| **ScanReports** | `GetReport`, `List`, `GetSummary` | `client.ScanReports` |
| **Attestation** | `GenerateLoi25` | `client.Attestation` |
| **Transactions** | `Create`, `List`, `Get`, `Delete`, `AddItem`, `UpdateItem`, `DeleteItem` | `client.Transactions` |
| **Templates** | `Create`, `List`, `Generate`, `ListDocuments` | `client.Templates` |
| **Batches** | `Create`, `List`, `Get`, `Commit`, `Cancel` | `client.Batches` |
| **Metadata** | `AddTags`, `GetTags`, `ListByTag`, `RemoveTag`, `SetMetadata`, `GetMetadata`, `DeleteMetadataKey` | `client.Metadata` |
| **Retention** | `Create`, `List`, `Update`, `Delete` | `client.Retention` |
| **LegalHolds** | `Create`, `List`, `Release`, `AddFiles`, `RemoveFile` | `client.LegalHolds` |
| **Permissions** | `Set`, `Get`, `Revoke` | `client.Permissions` |
| **Comments** | `Create`, `List`, `Update`, `Delete`, `GetReplies` | `client.Comments` |
| **Quota** | `Get` | `client.Quota` |
| **RateLimit** | `Get` | `client.RateLimit` |
| **UploadSessions** | `Create`, `UploadChunk`, `Complete`, `Cancel` | `client.UploadSessions` |
| **Jobs** | `Create`, `List`, `Get`, `Cancel` | `client.Jobs` |
| **ActivitySubscriptions** | `Subscribe`, `List`, `Unsubscribe` | `client.ActivitySubscriptions` |
| **Policies** | `GetIPPolicy`, `SetIPPolicy`, `GetMFAPolicy`, `SetMFAPolicy`, `GetEncryptionSalt`, `SetEncryptionSalt` | `client.Policies` |

## Files

```go
ctx := context.Background()

// List files (with optional filters)
files, err := client.Files.List(ctx, &cv.FileListOptions{
	FolderID: cv.String("folder-id"),
	Page:     1,
	Limit:    20,
})

// Get file metadata
file, err := client.Files.Get(ctx, "file-id")

// Upload a file
f, _ := os.Open("report.pdf")
defer f.Close()
result, err := client.Files.Upload(ctx, f, "report.pdf", cv.String("folder-id"))

// Download a file (returns io.ReadCloser -- caller must close)
reader, err := client.Files.Download(ctx, "file-id")
if err == nil {
	defer reader.Close()
	io.Copy(outFile, reader)
}

// Delete a file
err = client.Files.Delete(ctx, "file-id")
```

## Folders

```go
// List folders
folders, err := client.Folders.List(ctx, &cv.FolderListOptions{
	ParentID: cv.String("parent-id"),
})

// Get a folder
folder, err := client.Folders.Get(ctx, "folder-id")

// Create a folder
newFolder, err := client.Folders.Create(ctx, cv.CreateFolderRequest{
	Name:     "Reports",
	ParentID: cv.String("parent-folder-id"),
})

// Delete a folder
err = client.Folders.Delete(ctx, "folder-id")
```

## Share Links

```go
// List share links
links, err := client.ShareLinks.List(ctx, nil)

// Create a download share link (expires in 24 hours)
link, err := client.ShareLinks.Create(ctx, cv.CreateShareLinkRequest{
	FileID:    cv.String("file-id"),
	Type:      "download",
	ExpiresIn: 86400,
})
fmt.Printf("Share URL: %s\n", link.URL)

// Create a password-protected upload link
link, err = client.ShareLinks.Create(ctx, cv.CreateShareLinkRequest{
	FolderID: cv.String("folder-id"),
	Type:     "upload",
	Password: cv.String("s3cret"),
})

// Delete a share link
err = client.ShareLinks.Delete(ctx, "link-id")
```

## Electronic Signatures

```go
// Create a signature envelope
envelope, err := client.Signatures.Create(ctx, cv.CreateSignatureRequest{
	FileID:  "file-id",
	Subject: "Please sign this NDA",
	Signers: []cv.CreateSignatureSigner{
		{
			Email:     "signer@example.com",
			Name:      "Jane Doe",
			Role:      "signer",
			SignOrder: 1,
		},
	},
	ExpiryDays: 30,
})

// Check signature status
status, err := client.Signatures.GetStatus(ctx, envelope.ID)
fmt.Printf("Status: %s\n", status.Status)

// Download the signed document
reader, err := client.Signatures.DownloadSigned(ctx, envelope.ID)
if err == nil {
	defer reader.Close()
	// write to file...
}

// Revoke a pending envelope
err = client.Signatures.Revoke(ctx, envelope.ID)
```

## Webhooks

```go
// Register a webhook endpoint
resp, err := client.Webhooks.Register(ctx, cv.RegisterWebhookRequest{
	URL:    "https://your-app.com/webhooks/conformvault",
	Events: []string{"file.uploaded", "signature.completed"},
})
fmt.Printf("Secret (save this!): %s\n", resp.Secret)

// List all webhooks
hooks, err := client.Webhooks.List(ctx)

// Send a test event
err = client.Webhooks.Test(ctx, "webhook-id")

// Delete a webhook
err = client.Webhooks.Delete(ctx, "webhook-id")
```

### Webhook Signature Verification

Verify incoming webhook payloads in your HTTP handler using HMAC-SHA256:

```go
import cv "github.com/secuaas/conformvault-sdk-go"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	sigHeader := r.Header.Get("X-ConformVault-Signature")
	secret := "whsec_your_secret" // stored from Register response

	if !cv.VerifyWebhookSignature(payload, sigHeader, secret) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	// Signature is valid -- process the event
	fmt.Fprintf(w, "OK")
}
```

## Audit Logs

```go
entries, err := client.Audit.List(ctx, &cv.AuditListOptions{
	EventType: "file.uploaded",
	From:      "2025-01-01",
	To:        "2025-12-31",
	Page:      1,
	Limit:     50,
})
for _, e := range entries {
	fmt.Printf("[%s] %s on %s/%s\n", e.CreatedAt, e.Action, e.ResourceType, e.ResourceID)
}
```

## API Keys

```go
// List all keys
keys, err := client.Keys.List(ctx)

// Create a new key
newKey, err := client.Keys.Create(ctx, cv.CreateAPIKeyRequest{
	Name:        "CI/CD Key",
	Environment: "test",
	Scopes:      []string{"files:read", "files:write"},
})
fmt.Printf("Key (save this!): %s\n", newKey.Key)

// Get key metadata
key, err := client.Keys.Get(ctx, "key-id")

// Rotate a key (returns new key value)
rotated, err := client.Keys.Rotate(ctx, "key-id")
fmt.Printf("New key: %s\n", rotated.Key)

// Revoke a key
err = client.Keys.Revoke(ctx, "key-id")
```

## Bulk Operations

```go
// Bulk delete files
result, err := client.Bulk.Delete(ctx, []string{"file-1", "file-2", "file-3"})
fmt.Printf("Deleted %d/%d files\n", result.Succeeded, result.Processed)

// Bulk move files to a folder
result, err = client.Bulk.Move(ctx, []string{"file-1", "file-2"}, "target-folder-id")

// Bulk download as ZIP (returns io.ReadCloser)
zipReader, err := client.Bulk.Download(ctx, []string{"file-1", "file-2"})
if err == nil {
	defer zipReader.Close()
	io.Copy(zipFile, zipReader)
}
```

## File Versions

```go
// List all versions of a file
versions, err := client.Versions.List(ctx, "file-id")

// Get a specific version
version, err := client.Versions.Get(ctx, "file-id", "version-id")

// Restore an old version as current
err = client.Versions.Restore(ctx, "file-id", "version-id")

// Permanently delete a version
err = client.Versions.Delete(ctx, "file-id", "version-id")
```

## Search

```go
results, err := client.Search.Search(ctx, cv.SearchOptions{
	Query:    "quarterly report",
	Types:    "files,folders",
	Page:     1,
	PageSize: 20,
})
for _, r := range results.Data {
	fmt.Printf("[%s] %s -- %s\n", r.Type, r.Name, r.ID)
}
fmt.Printf("Total results: %d\n", results.Pagination.Total)
```

## Trash

```go
// List trashed files
trashed, err := client.Trash.List(ctx, &cv.TrashListOptions{Page: 1, Limit: 50})

// Restore a file from trash
err = client.Trash.Restore(ctx, "file-id")

// Permanently delete a single trashed file
err = client.Trash.Delete(ctx, "file-id")

// Empty the entire trash
result, err := client.Trash.Empty(ctx)
fmt.Printf("Permanently deleted %d files\n", result.FilesDeleted)
```

## Scan Reports

```go
// Get the scan report for a specific file
report, err := client.ScanReports.GetReport(ctx, "file-id")
fmt.Printf("Scan status: %s, engine: %s\n", report.ScanStatus, report.ScanEngine)

// List all scan reports
reports, err := client.ScanReports.List(ctx, &cv.ScanReportListOptions{
	Limit:  50,
	Offset: 0,
})

// Get aggregate scan summary for the organization
summary, err := client.ScanReports.GetSummary(ctx)
fmt.Printf("Total: %d, Clean: %d, Infected: %d\n",
	summary.TotalScans, summary.CleanCount, summary.InfectedCount)
```

## Compliance Attestation

```go
// Generate a Loi 25 compliance attestation PDF
reader, err := client.Attestation.GenerateLoi25(ctx)
if err == nil {
	defer reader.Close()
	out, _ := os.Create("attestation-loi25.pdf")
	defer out.Close()
	io.Copy(out, reader)
}
```

## Transaction Folders

```go
// Create a transaction folder
tx, err := client.Transactions.Create(ctx, cv.CreateTransactionRequest{
	Name:    "Real Estate Closing Q1",
	DueDate: cv.String("2025-06-30"),
})

// Add checklist items
item, err := client.Transactions.AddItem(ctx, tx.ID, cv.CreateTransactionItemRequest{
	Label:    "Signed purchase agreement",
	Required: cv.Bool(true),
})

// Update an item (e.g., attach a file)
_, err = client.Transactions.UpdateItem(ctx, tx.ID, item.ID, cv.UpdateTransactionItemRequest{
	Status: cv.String("completed"),
	FileID: cv.String("file-id"),
})

// List transaction folders
list, err := client.Transactions.List(ctx, 1, 20)

// Get a single transaction folder with items and progress
tx, err = client.Transactions.Get(ctx, "tx-id")
fmt.Printf("Progress: %d/%d completed\n", tx.Progress.Completed, tx.Progress.Total)

// Delete an item or the entire transaction
err = client.Transactions.DeleteItem(ctx, "tx-id", "item-id")
err = client.Transactions.Delete(ctx, "tx-id")
```

## Document Templates

```go
// Create a template
tmpl, err := client.Templates.Create(ctx, cv.CreateTemplateRequest{
	Name:        "Invoice Template",
	ContentType: "application/pdf",
	Fields:      []string{"client_name", "amount", "due_date"},
})

// Generate a PDF from the template (returns io.ReadCloser)
reader, err := client.Templates.Generate(ctx, tmpl.ID, cv.GenerateDocumentRequest{
	Data: map[string]string{
		"client_name": "Acme Corp",
		"amount":      "$5,000.00",
		"due_date":    "2025-04-15",
	},
})
if err == nil {
	defer reader.Close()
	// write to file...
}

// List generated documents for a template
docs, err := client.Templates.ListDocuments(ctx, tmpl.ID)
```

## Batch Uploads

```go
// Create a batch operation
batch, err := client.Batches.Create(ctx, cv.CreateBatchRequest{
	Type:     "upload",
	FolderID: cv.String("folder-id"),
	Items: []cv.CreateBatchItemDef{
		{Filename: "doc1.pdf", Size: 102400, MimeType: "application/pdf"},
		{Filename: "doc2.pdf", Size: 204800, MimeType: "application/pdf"},
	},
})

// Commit the batch to start processing
batch, err = client.Batches.Commit(ctx, batch.ID)

// Check batch status
batch, err = client.Batches.Get(ctx, batch.ID)
fmt.Printf("Status: %s (%d/%d completed)\n", batch.Status, batch.Completed, batch.Total)

// List all batch operations
list, err := client.Batches.List(ctx, 1, 20)

// Cancel a batch
err = client.Batches.Cancel(ctx, batch.ID)
```

## File Metadata & Tags

```go
// Add tags to a file
tags, err := client.Metadata.AddTags(ctx, "file-id", cv.AddTagsRequest{
	Tags: []string{"confidential", "legal", "q1-2025"},
})

// Get all tags for a file
tags, err = client.Metadata.GetTags(ctx, "file-id")

// List all files with a specific tag
files, err := client.Metadata.ListByTag(ctx, "confidential")

// Remove a tag
err = client.Metadata.RemoveTag(ctx, "file-id", "q1-2025")

// Set metadata key-value pairs
meta, err := client.Metadata.SetMetadata(ctx, "file-id", cv.SetMetadataRequest{
	Metadata: map[string]string{
		"department": "legal",
		"case_number": "2025-001",
	},
})

// Get metadata
meta, err = client.Metadata.GetMetadata(ctx, "file-id")

// Delete a metadata key
err = client.Metadata.DeleteMetadataKey(ctx, "file-id", "case_number")
```

## Retention Policies

```go
// Create a retention policy
policy, err := client.Retention.Create(ctx, cv.CreateRetentionPolicyRequest{
	Name:          "7-Year Legal Hold",
	RetentionDays: 2555,
	AutoDelete:    false,
})

// List all policies
policies, err := client.Retention.List(ctx)

// Update a policy
updated, err := client.Retention.Update(ctx, policy.ID, cv.UpdateRetentionPolicyRequest{
	AutoDelete: cv.Bool(true),
})

// Delete a policy
err = client.Retention.Delete(ctx, policy.ID)
```

## Legal Holds

```go
// Create a legal hold
hold, err := client.LegalHolds.Create(ctx, cv.CreateLegalHoldRequest{
	Name:        "Case 2025-001",
	Description: "Litigation hold for pending lawsuit",
})

// Add files to the hold
files, err := client.LegalHolds.AddFiles(ctx, hold.ID, cv.AddLegalHoldFilesRequest{
	FileIDs: []string{"file-1", "file-2", "file-3"},
})

// List all holds
holds, err := client.LegalHolds.List(ctx)

// Release a hold
released, err := client.LegalHolds.Release(ctx, hold.ID)

// Remove a file from a hold
err = client.LegalHolds.RemoveFile(ctx, hold.ID, "file-1")
```

## Folder Permissions

```go
// Set permissions on a folder
perm, err := client.Permissions.Set(ctx, "folder-id", cv.SetFolderPermissionRequest{
	UserID:     "user-id",
	Permission: "write",
})

// Get folder permissions
perms, err := client.Permissions.Get(ctx, "folder-id")

// Revoke a user's permission
err = client.Permissions.Revoke(ctx, "folder-id", "user-id")
```

## Comments

```go
// Add a comment to a file
comment, err := client.Comments.Create(ctx, "file-id", cv.CreateCommentRequest{
	Content: "Please review section 3 before signing.",
})

// Reply to a comment
reply, err := client.Comments.Create(ctx, "file-id", cv.CreateCommentRequest{
	Content:  "Section 3 looks good to me.",
	ParentID: &comment.ID,
})

// List comments on a file
comments, err := client.Comments.List(ctx, "file-id")

// Get replies to a comment
replies, err := client.Comments.GetReplies(ctx, "file-id", comment.ID)

// Update a comment
updated, err := client.Comments.Update(ctx, "file-id", comment.ID, cv.UpdateCommentRequest{
	Content: "Updated: Please review sections 3 and 4.",
})

// Delete a comment
err = client.Comments.Delete(ctx, "file-id", comment.ID)
```

## Quota & Rate Limits

```go
// Get storage quota
quota, err := client.Quota.Get(ctx)
fmt.Printf("Used: %d / %d bytes (%d files)\n",
	quota.UsedBytes, quota.TotalBytes, quota.FileCount)

// Get rate limit status
rl, err := client.RateLimit.Get(ctx)
fmt.Printf("Remaining: %d/%d requests (resets at %s)\n",
	rl.RequestsRemaining, rl.RequestsPerMinute, rl.ResetAt)
```

## Upload Sessions (Chunked Uploads)

```go
// Create an upload session for a large file
session, err := client.UploadSessions.Create(ctx, cv.CreateUploadSessionRequest{
	Filename:    "large-backup.zip",
	TotalSize:   5368709120, // 5 GB
	ContentType: "application/zip",
})

// Upload chunks
for i := 0; i < session.TotalChunks; i++ {
	chunk := getChunkReader(i) // your chunk data
	err = client.UploadSessions.UploadChunk(ctx, session.ID, i, chunk)
}

// Complete the session
file, err := client.UploadSessions.Complete(ctx, session.ID)
fmt.Printf("Upload complete: %s\n", file.ID)

// Or cancel if needed
err = client.UploadSessions.Cancel(ctx, session.ID)
```

## Background Jobs

```go
// Create a background job
job, err := client.Jobs.Create(ctx, cv.CreateJobRequest{
	Type:   "export",
	Params: map[string]any{"folder_id": "folder-id", "format": "zip"},
})

// Poll job status
job, err = client.Jobs.Get(ctx, job.ID)
fmt.Printf("Job %s: %s (%d%% complete)\n", job.ID, job.Status, job.Progress)

// List all jobs
jobs, err := client.Jobs.List(ctx)

// Cancel a job
err = client.Jobs.Cancel(ctx, job.ID)
```

## Activity Subscriptions

```go
// Subscribe to activity events
sub, err := client.ActivitySubscriptions.Subscribe(ctx, cv.CreateActivitySubscriptionRequest{
	EventTypes:  []string{"file.uploaded", "file.deleted"},
	CallbackURL: "https://your-app.com/activity",
})

// List subscriptions
subs, err := client.ActivitySubscriptions.List(ctx)

// Unsubscribe
err = client.ActivitySubscriptions.Unsubscribe(ctx, sub.ID)
```

## Security Policies

```go
// Get and set IP restriction policy
ipPolicy, err := client.Policies.GetIPPolicy(ctx)
ipPolicy, err = client.Policies.SetIPPolicy(ctx, cv.SetIPPolicyRequest{
	Enabled:    true,
	AllowedIPs: []string{"203.0.113.0/24", "198.51.100.1"},
})

// Get and set MFA policy
mfaPolicy, err := client.Policies.GetMFAPolicy(ctx)
mfaPolicy, err = client.Policies.SetMFAPolicy(ctx, cv.SetMFAPolicyRequest{
	Enabled:     true,
	RequiredFor: []string{"file.delete", "settings.update"},
})

// Get and set encryption salt
salt, err := client.Policies.GetEncryptionSalt(ctx)
salt, err = client.Policies.SetEncryptionSalt(ctx, cv.SetEncryptionSaltRequest{
	Salt: "your-base64-encoded-salt",
})
```

## Error Handling

All API errors are returned as typed errors that can be inspected:

```go
file, err := client.Files.Get(ctx, "nonexistent-id")
if err != nil {
	// Check for specific error types
	if cv.IsNotFound(err) {
		fmt.Println("File not found")
		return
	}

	if cv.IsRateLimited(err) {
		rlErr := err.(*cv.RateLimitError)
		fmt.Printf("Rate limited -- retry after %s\n", rlErr.RetryAfter)
		return
	}

	// Generic API error
	if apiErr, ok := err.(*cv.APIError); ok {
		fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
		return
	}

	// Network or other error
	fmt.Printf("Unexpected error: %v\n", err)
}
```

### Error Types

| Type | Description | Helper |
|------|-------------|--------|
| `*APIError` | Any HTTP 4xx/5xx response from the API | -- |
| `*RateLimitError` | HTTP 429 with `RetryAfter` duration | `IsRateLimited(err)` |
| -- | HTTP 404 (a subset of `*APIError`) | `IsNotFound(err)` |

### `APIError` Fields

```go
type APIError struct {
	StatusCode int    // HTTP status code
	Message    string // Error message from the API
}
```

### `RateLimitError` Fields

```go
type RateLimitError struct {
	APIError               // Embeds APIError
	RetryAfter time.Duration // Parsed from Retry-After header
}
```

## Helper Functions

The SDK provides pointer-helper functions for optional fields:

```go
// Built-in helpers
cv.String("value")  // returns *string
cv.Bool(true)       // returns *bool
cv.Int(42)          // returns *int

// Usage
client.Files.Upload(ctx, reader, "file.pdf", cv.String("folder-id"))
client.ShareLinks.Create(ctx, cv.CreateShareLinkRequest{
	FileID:   cv.String("file-id"),
	Type:     "download",
	Password: cv.String("s3cret"),
})
```

## Thread Safety

The `Client` and all service objects are safe for concurrent use across goroutines. The underlying `*http.Client` handles connection pooling.

## License

MIT -- see [LICENSE](LICENSE).
