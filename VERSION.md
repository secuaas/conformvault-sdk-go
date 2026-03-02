# Historique des Versions - ConformVault Go SDK

## Version Actuelle
**0.4.0** - 2026-03-02

---

## Versions

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

### 0.2.0 - 2026-02-27
**Type:** Minor - V2-4 SDK expansion: bulk, versions, search, trash services

### Ajouté
- `BulkService`, `VersionsService`, `SearchService`, `TrashService`

---

### 0.1.0 - 2026-02-26
**Commit:** `f67c1f2`
**Type:** Initial release — 7 services (files, folders, sharelinks, signatures, webhooks, audit, keys)
