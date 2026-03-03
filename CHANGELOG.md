# Changelog

All notable changes to the ConformVault Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
