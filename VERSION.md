# Historique des Versions - ConformVault Go SDK

## Version Actuelle
**0.5.0** - 2026-03-04

---

## Versions

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
