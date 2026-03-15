# Historique des Versions - ConformVault Go SDK

## Version Actuelle
**0.6.0** - 2026-03-15

---

## Versions

### 0.6.0 - 2026-03-15
**Commit:** `pending`
**Type:** Minor - Add 8 v3 feature services: SecretVault, ExpectedFiles, SpaceMessaging, Retention Exceptions/Approvals, Temporary Permissions, Signature Delegation, MSP Dashboard, Import Wizard

### Ajouté
- **`secret_vault.go`**: `SecretVaultService` — Create, List, Get, Delete (4 methods)
- **`expected_files.go`**: `ExpectedFilesService` — Create, List, Update, Delete, GetProgress (5 methods)
- **`space_messaging.go`**: `SpaceMessagingService` — SendMessage, ListMessages, GetReplies, MarkRead, DeleteMessage (5 methods)
- **`msp_dashboard.go`**: `MSPDashboardService` — GetDashboard, ListClients, GetClientUsage (3 methods)
- **`imports.go`**: `ImportsService` — CreateConnection, ListConnections, DeleteConnection, StartImport, ListJobs, GetJob, CancelJob (7 methods)
- **`retention.go`**: Added CreateException, ListExceptions, DeleteException, ListPendingApprovals, DecideApproval (5 methods)
- **`permissions.go`**: Added SetWithExpiry for temporary access (1 method)
- **`signatures.go`**: Added Delegate for signature delegation (1 method)
- 20+ new types in `types.go` (Secret, ExpectedFile, SpaceMessage, RetentionException, RetentionApproval, DelegateSignRequest, MSPDashboardSummary, MSPClientMetrics, ImportConnection, ImportJob, etc.)
- `ExpiresAt` field added to `FolderPermission` type
- 34 new unit tests (51 → 85 total)
- Total services: 29 → 34

### Tests effectués
- ✅ `go build ./...` — success
- ✅ `go test ./...` — 85 tests passing

---

### 0.5.4 - 2026-03-04
**Commit:** `pending`
**Type:** Patch - Add pagination support to Webhooks.List, Keys.List + PaginationInfo type

### Ajouté
- `PaginationInfo` type in `types.go` — reusable pagination metadata struct (page, limit, total)
- `WebhookListOptions` type — query parameters for paginated webhook listing
- `KeyListOptions` type — query parameters for paginated API key listing

### Modifié
- `WebhooksService.List` now accepts `*WebhookListOptions` for optional `page` and `limit` query params
- `KeysService.List` now accepts `*KeyListOptions` for optional `page` and `limit` query params
- Updated tests to pass `nil` opts to maintain backward-compatible behavior

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.5.3 - 2026-03-04
**Commit:** `pending`
**Type:** Patch - Add 3 new signature methods (AnalyzePDF, PreviewPDF, GetEmbeddedSignLink)

### Ajouté
- `SignaturesService.AnalyzePDF` — Analyze a PDF for signature field placement (POST /signatures/analyze)
- `SignaturesService.PreviewPDF` — Stream a decrypted PDF for signature placement preview (GET /signatures/preview-pdf)
- `SignaturesService.GetEmbeddedSignLink` — Retrieve an embedded signing link for a signer (GET /signatures/{id}/embed-sign)
- New types: `PDFPageDimension`, `PDFFieldSuggestion`, `PDFAnalysisResult`, `AnalyzePDFRequest`, `EmbeddedSignLinkResponse`

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.5.2 - 2026-03-04
**Commit:** `pending`
**Type:** Patch - Add unit tests (51 tests via go test)

### Ajouté
- 51 unit tests using go test — covers all services with httptest mock server
- Test infrastructure with net/http/httptest

### Tests effectués
- ✅ `go build ./...` — success
- ✅ `go test ./...` — 51 tests passing

---

### 0.5.1 - 2026-03-04
**Commit:** `pending`
**Type:** Patch - Fix 9 route mismatches + add 6 missing routes (backend alignment)

### Corrigé
- **Policies**: `/policies/ip` → `/ip-policy`, `/policies/mfa` → `/mfa-policy`, `/policies/encryption-salt` → `/encryption/salt`
- **Activity Subscriptions**: `/activity/subscriptions` → `/activity-subscriptions`
- **Comments**: Rewritten from nested `/files/{id}/comments` to flat `/comments` routes; `file_id` moved to request body/query param; Update changed from PATCH to PUT
- **Jobs Cancel**: `POST /jobs/{id}/cancel` → `DELETE /jobs/{id}`
- **Batches Cancel**: `POST /batches/{id}/cancel` → `DELETE /batches/{id}`
- **Retention Update**: PATCH → PUT

### Ajouté
- **`bandwidth.go`**: `BandwidthService` — GetSummary (`GET /bandwidth`), GetDaily (`GET /bandwidth/daily`)
- **`data_export.go`**: `DataExportService` — Export (`GET /users/{id}/export`) for GDPR/Loi 25
- **`keys.go`**: InstantRevoke (`POST /api-keys/{id}/revoke`), GetRevocationStatus (`GET /api-keys/{id}/revocation-status`)
- **`batches.go`**: UploadFile (`PUT /batches/{id}/files/{index}`)
- New types: `BandwidthSummary`, `DailyBandwidthStats`, `KeyRevocationStatus`, `UserDataExport`
- `FileID` field added to `CreateCommentRequest`
- Total services: 27 → 29

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.5.0 - 2026-03-04
**Commit:** `pending`
**Type:** Minor - SDK v0.5.0: 57 new methods across 11 new services (~85% API coverage)

### Ajouté
- **Webhooks**: `ListDeliveries`, `GetDelivery`, `ReplayDelivery`, `ReEnable` (4 methods)
- **Audit**: `Search`, `Export` (raw stream), `GetStats`, `GetAnomalies` (4 methods)
- **Files**: `GetThumbnail` (raw stream), `GetScanReport` (2 methods)
- **`metadata.go`**: `MetadataService` — AddTags, RemoveTag, GetTags, ListByTag, SetMetadata, GetMetadata, DeleteMetadataKey (7 methods)
- **`retention.go`**: `RetentionService` — Create, List, Get, Update, Delete (5 methods)
- **`legal_holds.go`**: `LegalHoldsService` — Create, List, Get, Release, AddFiles, RemoveFile (6 methods)
- **`permissions.go`**: `PermissionsService` — Set, Get, Revoke (3 methods)
- **`comments.go`**: `CommentsService` — Create, List, Get, Update, Delete, GetReplies (6 methods)
- **`quota.go`**: `QuotaService` — Get (1 method) + `RateLimitService` — Get (1 method)
- **`upload_sessions.go`**: `UploadSessionsService` — Create, UploadChunk, GetStatus, Complete, Cancel (5 methods)
- **`jobs.go`**: `JobsService` — Create, List, Get, Cancel (4 methods)
- **`activity_subscriptions.go`**: `ActivitySubscriptionsService` — Subscribe, List, Unsubscribe (3 methods)
- **`policies.go`**: `PoliciesService` — GetIPPolicy, SetIPPolicy, GetMFAPolicy, SetMFAPolicy, GetEncryptionSalt, SetEncryptionSalt (6 methods)
- 30+ new types in `types.go`
- Version bumped to 0.5.0
- Total services: 16 → 27

### Tests effectués
- ✅ `go build ./...` — success
- ✅ `go vet ./...` — success

---

### 0.4.1 - 2026-03-03
**Commit:** `50e07df`
**Type:** Patch - Add LICENSE (MIT) and CHANGELOG.md

### Ajouté
- `LICENSE` — MIT license file (was missing, referenced by README badge)
- `CHANGELOG.md` — Keep a Changelog format, extracted from VERSION.md history

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.4.0 - 2026-03-02
**Commit:** `d821e0a`
**Type:** Minor - Transactions, Templates, Batches services

### Ajouté
- **`transactions.go`**: `TransactionsService` — Create, List, Get, Update, Delete, AddItem, UpdateItem, DeleteItem (8 methods)
- **`templates.go`**: `TemplatesService` — Create, List, Get, Update, Delete, Generate (binary PDF via `doRaw`), ListDocuments (7 methods)
- **`batches.go`**: `BatchesService` — Create, List, Get, Commit, Cancel (5 methods)
- New types: `TransactionFolder`, `TransactionFolderItem`, `TransactionProgress`, `DocumentTemplate`, `GeneratedDocument`, `BatchOperation`, `BatchOperationItem` + request/response types
- Services registered as `client.Transactions`, `client.Templates`, `client.Batches`
- Total services: 13 → 16

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.3.0 - 2026-02-27
**Commit:** `507eaf3`
**Type:** Minor - ScanReports and Attestation services

### Ajouté
- **`scan_reports.go`**: `ScanReportsService` — GetReport, List, GetSummary
- **`attestation.go`**: `AttestationService` — GenerateLoi25 (PDF stream)
- New types: `FileScanReport`, `FileScanSummary`, `ScanReportListOptions`, `ScanReportListResponse`
- Services registered in Client struct

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.2.0 - 2026-02-27 (backfill committed 2026-03-02)
**Commit:** `df74b64`
**Type:** Minor - V2-4 SDK expansion: bulk, versions, search, trash services

### Ajouté
- **`bulk.go`**: `BulkService` — Delete, Move, Download (ZIP stream)
- **`versions.go`**: `VersionsService` — List, Get, Restore, Delete
- **`search.go`**: `SearchService` — Search (with query params: q, types, folder_id, page, page_size)
- **`trash.go`**: `TrashService` — List, Restore, Delete, Empty, GetStats

### Tests effectués
- ✅ `go build ./...` — success

---

### 0.1.0 - 2026-02-26
**Commit:** `f67c1f2`
**Type:** Initial release — 7 services (files, folders, sharelinks, signatures, webhooks, audit, keys)
