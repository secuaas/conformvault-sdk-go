# Changelog

All notable changes to the ConformVault Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.2] - 2026-03-04

### Added
- 51 unit tests using `go test` — covers all 29 services with `httptest` mock server
- Test infrastructure with `net/http/httptest`

## [0.5.1] - 2026-03-04

### Fixed
- **Policies**: paths corrected from `/policies/ip` to `/ip-policy`, `/policies/mfa` to `/mfa-policy`, `/policies/encryption-salt` to `/encryption/salt`
- **Activity Subscriptions**: path corrected from `/activity/subscriptions` to `/activity-subscriptions`
- **Comments**: rewritten from nested `/files/{id}/comments` to flat `/comments` routes; Update method changed from PATCH to PUT
- **Jobs Cancel**: changed from `POST /jobs/{id}/cancel` to `DELETE /jobs/{id}`
- **Batches Cancel**: changed from `POST /batches/{id}/cancel` to `DELETE /batches/{id}`
- **Retention Update**: changed from PATCH to PUT

### Added
- `BandwidthService` — GetSummary, GetDaily (2 methods)
- `DataExportService` — Export for GDPR/Loi 25 (1 method)
- `KeysService`: InstantRevoke, GetRevocationStatus (2 new methods via `/api-keys/` path)
- `BatchesService`: UploadFile (1 new method)
- New types: `BandwidthSummary`, `DailyBandwidthStats`, `KeyRevocationStatus`, `UserDataExport`
- `FileID` field added to `CreateCommentRequest`
- Total services: 27 → 29

## [0.5.0] - 2026-03-04

### Added
- `MetadataService` -- AddTags, RemoveTag, GetTags, ListByTag, SetMetadata, GetMetadata, DeleteMetadataKey (7 methods)
- `RetentionService` -- Create, List, Get, Update, Delete (5 methods)
- `LegalHoldsService` -- Create, List, Get, Release, AddFiles, RemoveFile (6 methods)
- `PermissionsService` -- Set, Get, Revoke (3 methods)
- `CommentsService` -- Create, List, Get, Update, Delete, GetReplies (6 methods)
- `QuotaService` -- Get (1 method)
- `RateLimitService` -- Get (1 method)
- `UploadSessionsService` -- Create, UploadChunk, GetStatus, Complete, Cancel (5 methods)
- `JobsService` -- Create, List, Get, Cancel (4 methods)
- `ActivitySubscriptionsService` -- Subscribe, List, Unsubscribe (3 methods)
- `PoliciesService` -- GetIPPolicy, SetIPPolicy, GetMFAPolicy, SetMFAPolicy, GetEncryptionSalt, SetEncryptionSalt (6 methods)
- `Webhooks`: ListDeliveries, GetDelivery, ReplayDelivery, ReEnable (4 new methods)
- `Audit`: Search, Export (raw stream), GetStats, GetAnomalies (4 new methods)
- `Files`: GetThumbnail (raw stream), GetScanReport (2 new methods)
- Pointer helper functions: `String()`, `Bool()`, `Int()`
- 30+ new types in `types.go`
- Total services: 16 -> 27

### Fixed
- URL-encode all query parameter values using `net/url.QueryEscape` (prevents injection via crafted filter values)
- Fix double-read of HTTP error response body in `do()` method; also cap error body reads to 1 MB

### Changed
- CONTRIBUTING.md and SECURITY.md added for open-source readiness

## [0.4.0] - 2026-03-02

### Added
- `TransactionsService` — Create, List, Get, Update, Delete, AddItem, UpdateItem, DeleteItem (8 methods)
- `TemplatesService` — Create, List, Get, Update, Delete, Generate (binary PDF via `doRaw`), ListDocuments (7 methods)
- `BatchesService` — Create, List, Get, Commit, Cancel (5 methods)
- New types: `TransactionFolder`, `TransactionFolderItem`, `TransactionProgress`, `DocumentTemplate`, `GeneratedDocument`, `BatchOperation`, `BatchOperationItem`
- Total services: 13 → 16

## [0.3.0] - 2026-02-27

### Added
- `ScanReportsService` — GetReport, List, GetSummary
- `AttestationService` — GenerateLoi25 (PDF stream)
- New types: `FileScanReport`, `FileScanSummary`, `ScanReportListOptions`, `ScanReportListResponse`

## [0.2.0] - 2026-02-27

### Added
- `BulkService` — Delete, Move, Download (ZIP stream)
- `VersionsService` — List, Get, Restore, Delete
- `SearchService` — Search (with query params: q, types, folder_id, page, page_size)
- `TrashService` — List, Restore, Delete, Empty, GetStats

## [0.1.0] - 2026-02-26

### Added
- Initial release with 7 services: files, folders, sharelinks, signatures, webhooks, audit, keys
- Full Developer API coverage for core ConformVault operations
