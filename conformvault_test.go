package conformvault

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helper: create a mock server that validates method, path, and auth header,
// then returns the given JSON payload.
// ---------------------------------------------------------------------------

type requestLog struct {
	Method string
	Path   string
	Query  string
	Auth   string
	UA     string
	CT     string
	Body   string
}

// newMockServer creates a test server that records every request into *logs
// and dispatches to handler based on "METHOD path".
func newMockServer(t *testing.T, routes map[string]http.HandlerFunc) (*httptest.Server, *[]requestLog) {
	t.Helper()
	var logs []requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read body
		var body string
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			body = string(b)
		}
		logs = append(logs, requestLog{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Auth:   r.Header.Get("Authorization"),
			UA:     r.Header.Get("User-Agent"),
			CT:     r.Header.Get("Content-Type"),
			Body:   body,
		})

		key := r.Method + " " + r.URL.Path
		if h, ok := routes[key]; ok {
			h(w, r)
			return
		}
		// fallback: try without trailing slash
		t.Errorf("unhandled route: %s %s (query: %s)", r.Method, r.URL.Path, r.URL.RawQuery)
		w.WriteHeader(http.StatusNotFound)
	}))
	return srv, &logs
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ---------------------------------------------------------------------------
// Client creation tests
// ---------------------------------------------------------------------------

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("test-key")
	if c.baseURL != DefaultBaseURL {
		t.Errorf("expected default base URL %q, got %q", DefaultBaseURL, c.baseURL)
	}
	if c.apiKey != "test-key" {
		t.Errorf("expected api key 'test-key', got %q", c.apiKey)
	}
	if c.Files == nil {
		t.Error("expected Files service to be initialized")
	}
	if c.Folders == nil {
		t.Error("expected Folders service to be initialized")
	}
	if c.Signatures == nil {
		t.Error("expected Signatures service to be initialized")
	}
	if c.ShareLinks == nil {
		t.Error("expected ShareLinks service to be initialized")
	}
	if c.Webhooks == nil {
		t.Error("expected Webhooks service to be initialized")
	}
	if c.Audit == nil {
		t.Error("expected Audit service to be initialized")
	}
	if c.Keys == nil {
		t.Error("expected Keys service to be initialized")
	}
	if c.Comments == nil {
		t.Error("expected Comments service to be initialized")
	}
	if c.Jobs == nil {
		t.Error("expected Jobs service to be initialized")
	}
	if c.Bandwidth == nil {
		t.Error("expected Bandwidth service to be initialized")
	}
	if c.DataExport == nil {
		t.Error("expected DataExport service to be initialized")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("key", WithBaseURL("https://custom.api.com"), WithHTTPClient(hc))
	if c.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom base URL, got %q", c.baseURL)
	}
	if c.httpClient != hc {
		t.Error("expected custom HTTP client")
	}
}

func TestNewClient_AllServicesNonNil(t *testing.T) {
	c := NewClient("k")
	services := map[string]any{
		"Files":                 c.Files,
		"Folders":               c.Folders,
		"ShareLinks":            c.ShareLinks,
		"Signatures":            c.Signatures,
		"Webhooks":              c.Webhooks,
		"Audit":                 c.Audit,
		"Keys":                  c.Keys,
		"Bulk":                  c.Bulk,
		"Versions":              c.Versions,
		"Search":                c.Search,
		"Trash":                 c.Trash,
		"ScanReports":           c.ScanReports,
		"Attestation":           c.Attestation,
		"Transactions":          c.Transactions,
		"Templates":             c.Templates,
		"Batches":               c.Batches,
		"Metadata":              c.Metadata,
		"Retention":             c.Retention,
		"LegalHolds":            c.LegalHolds,
		"Permissions":           c.Permissions,
		"Comments":              c.Comments,
		"Quota":                 c.Quota,
		"RateLimit":             c.RateLimit,
		"UploadSessions":        c.UploadSessions,
		"Jobs":                  c.Jobs,
		"ActivitySubscriptions": c.ActivitySubscriptions,
		"Policies":              c.Policies,
		"Bandwidth":             c.Bandwidth,
		"DataExport":            c.DataExport,
	}
	for name, svc := range services {
		if svc == nil {
			t.Errorf("expected %s service to be initialized, got nil", name)
		}
	}
}

// ---------------------------------------------------------------------------
// Auth header validation (shared across all requests)
// ---------------------------------------------------------------------------

func TestAuthHeader(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[File]{Data: []File{}})
		},
	})
	defer srv.Close()

	c := NewClient("cvk_live_abc123", WithBaseURL(srv.URL))
	_, _ = c.Files.List(context.Background(), nil)

	if len(*logs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(*logs))
	}
	if (*logs)[0].Auth != "Bearer cvk_live_abc123" {
		t.Errorf("expected Bearer auth, got %q", (*logs)[0].Auth)
	}
	if (*logs)[0].UA != userAgent {
		t.Errorf("expected User-Agent %q, got %q", userAgent, (*logs)[0].UA)
	}
}

// ---------------------------------------------------------------------------
// Files service
// ---------------------------------------------------------------------------

func TestFilesList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[File]{
				Data: []File{
					{ID: "f1", Name: "test.pdf", Size: 1024},
					{ID: "f2", Name: "doc.txt", Size: 512},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	files, err := c.Files.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].ID != "f1" {
		t.Errorf("expected first file ID 'f1', got %q", files[0].ID)
	}
	if files[1].Name != "doc.txt" {
		t.Errorf("expected second file name 'doc.txt', got %q", files[1].Name)
	}
	if (*logs)[0].Method != "GET" {
		t.Errorf("expected GET, got %s", (*logs)[0].Method)
	}
}

func TestFilesListWithOptions(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[File]{Data: []File{}})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	folderID := "folder-abc"
	_, err := c.Files.List(context.Background(), &FileListOptions{
		FolderID: &folderID,
		Page:     2,
		Limit:    10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q := (*logs)[0].Query
	if !strings.Contains(q, "folder_id=folder-abc") {
		t.Errorf("expected folder_id in query, got %q", q)
	}
	if !strings.Contains(q, "page=2") {
		t.Errorf("expected page=2 in query, got %q", q)
	}
	if !strings.Contains(q, "limit=10") {
		t.Errorf("expected limit=10 in query, got %q", q)
	}
}

func TestFilesGet(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files/file-123": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[File]{
				Data: File{ID: "file-123", Name: "report.pdf", Size: 2048},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	file, err := c.Files.Get(context.Background(), "file-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.ID != "file-123" {
		t.Errorf("expected file ID 'file-123', got %q", file.ID)
	}
	if file.Name != "report.pdf" {
		t.Errorf("expected file name 'report.pdf', got %q", file.Name)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/files/file-123" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestFilesDelete(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /files/file-456": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Files.Delete(context.Background(), "file-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/files/file-456" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Folders service
// ---------------------------------------------------------------------------

func TestFoldersList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /folders": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[Folder]{
				Data: []Folder{
					{ID: "fld-1", Name: "Documents", Path: "/Documents"},
					{ID: "fld-2", Name: "Images", Path: "/Images"},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	folders, err := c.Folders.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(folders) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(folders))
	}
	if folders[0].Name != "Documents" {
		t.Errorf("expected 'Documents', got %q", folders[0].Name)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/folders" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestFoldersCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /folders": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[Folder]{
				Data: Folder{ID: "fld-new", Name: "Projects", Path: "/Projects"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	folder, err := c.Folders.Create(context.Background(), CreateFolderRequest{
		Name: "Projects",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if folder.ID != "fld-new" {
		t.Errorf("expected folder ID 'fld-new', got %q", folder.ID)
	}
	if folder.Name != "Projects" {
		t.Errorf("expected folder name 'Projects', got %q", folder.Name)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/folders" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
	if (*logs)[0].CT != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", (*logs)[0].CT)
	}
}

func TestFoldersGet(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /folders/fld-1": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[Folder]{
				Data: Folder{ID: "fld-1", Name: "Documents", Path: "/Documents"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	folder, err := c.Folders.Get(context.Background(), "fld-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if folder.ID != "fld-1" {
		t.Errorf("expected folder ID 'fld-1', got %q", folder.ID)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/folders/fld-1" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestFoldersDelete(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /folders/fld-99": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Folders.Delete(context.Background(), "fld-99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/folders/fld-99" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// ShareLinks service
// ---------------------------------------------------------------------------

func TestShareLinksList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /sharelinks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[ShareLink]{
				Data: []ShareLink{
					{ID: "sl-1", Token: "abc", URL: "https://vault.test/s/abc", Type: "download", IsActive: true},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	links, err := c.ShareLinks.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 share link, got %d", len(links))
	}
	if links[0].Type != "download" {
		t.Errorf("expected type 'download', got %q", links[0].Type)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/sharelinks" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestShareLinksCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /sharelinks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[ShareLink]{
				Data: ShareLink{ID: "sl-new", Token: "xyz", URL: "https://vault.test/s/xyz", Type: "upload"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	fileID := "file-abc"
	link, err := c.ShareLinks.Create(context.Background(), CreateShareLinkRequest{
		FileID: &fileID,
		Type:   "upload",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.ID != "sl-new" {
		t.Errorf("expected share link ID 'sl-new', got %q", link.ID)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/sharelinks" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestShareLinksDelete(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /sharelinks/sl-del": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.ShareLinks.Delete(context.Background(), "sl-del")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/sharelinks/sl-del" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Signatures service
// ---------------------------------------------------------------------------

func TestSignaturesCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /signatures": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[SignatureEnvelope]{
				Data: SignatureEnvelope{
					ID:      "env-123",
					Status:  "sent",
					Subject: "NDA Agreement",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	env, err := c.Signatures.Create(context.Background(), CreateSignatureRequest{
		FileID:  "file-abc",
		Subject: "NDA Agreement",
		Signers: []CreateSignatureSigner{
			{Email: "signer@example.com", Name: "John Doe"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.ID != "env-123" {
		t.Errorf("expected envelope ID 'env-123', got %q", env.ID)
	}
	if env.Status != "sent" {
		t.Errorf("expected status 'sent', got %q", env.Status)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/signatures" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestSignaturesGetStatus(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /signatures/env-123": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[SignatureEnvelope]{
				Data: SignatureEnvelope{ID: "env-123", Status: "completed", Subject: "NDA"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	env, err := c.Signatures.GetStatus(context.Background(), "env-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Status != "completed" {
		t.Errorf("expected status 'completed', got %q", env.Status)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/signatures/env-123" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestSignaturesRevoke(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /signatures/env-456/revoke": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Signatures.Revoke(context.Background(), "env-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/signatures/env-456/revoke" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Webhooks service
// ---------------------------------------------------------------------------

func TestWebhooksList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /webhooks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[WebhookEndpoint]{
				Data: []WebhookEndpoint{
					{ID: "wh-1", URL: "https://example.com/hook", Events: []string{"file.uploaded"}, IsActive: true},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	hooks, err := c.Webhooks.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(hooks))
	}
	if hooks[0].ID != "wh-1" {
		t.Errorf("expected webhook ID 'wh-1', got %q", hooks[0].ID)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/webhooks" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestWebhooksRegister(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /webhooks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, struct {
				Data    RegisterWebhookResponse `json:"data"`
				Message string                  `json:"message"`
			}{
				Data: RegisterWebhookResponse{
					WebhookEndpoint: WebhookEndpoint{
						ID:     "wh-new",
						URL:    "https://example.com/hook",
						Events: []string{"file.uploaded", "file.deleted"},
					},
					Secret: "whsec_test_secret",
				},
				Message: "Webhook registered",
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	resp, err := c.Webhooks.Register(context.Background(), RegisterWebhookRequest{
		URL:    "https://example.com/hook",
		Events: []string{"file.uploaded", "file.deleted"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "wh-new" {
		t.Errorf("expected webhook ID 'wh-new', got %q", resp.ID)
	}
	if resp.Secret != "whsec_test_secret" {
		t.Errorf("expected secret 'whsec_test_secret', got %q", resp.Secret)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/webhooks" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestWebhooksTest(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /webhooks/wh-1/test": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Webhooks.Test(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/webhooks/wh-1/test" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestWebhooksListDeliveries(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /webhooks/wh-1/deliveries": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[WebhookDelivery]{
				Data: []WebhookDelivery{
					{ID: "del-1", WebhookID: "wh-1", EventType: "file.uploaded", Status: "success", HTTPStatus: 200},
					{ID: "del-2", WebhookID: "wh-1", EventType: "file.deleted", Status: "failed", HTTPStatus: 500},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	deliveries, err := c.Webhooks.ListDeliveries(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deliveries) != 2 {
		t.Fatalf("expected 2 deliveries, got %d", len(deliveries))
	}
	if deliveries[0].Status != "success" {
		t.Errorf("expected first delivery status 'success', got %q", deliveries[0].Status)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/webhooks/wh-1/deliveries" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Audit service
// ---------------------------------------------------------------------------

func TestAuditList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /audit": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[AuditEntry]{
				Data: []AuditEntry{
					{ID: "aud-1", Action: "file.uploaded", ResourceType: "file", ResourceID: "f1"},
					{ID: "aud-2", Action: "file.deleted", ResourceType: "file", ResourceID: "f2"},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	entries, err := c.Audit.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 audit entries, got %d", len(entries))
	}
	if entries[0].Action != "file.uploaded" {
		t.Errorf("expected action 'file.uploaded', got %q", entries[0].Action)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/audit" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestAuditListWithFilters(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /audit": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[AuditEntry]{Data: []AuditEntry{}})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := c.Audit.List(context.Background(), &AuditListOptions{
		EventType: "file.uploaded",
		From:      "2026-01-01",
		To:        "2026-03-01",
		Page:      1,
		Limit:     50,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q := (*logs)[0].Query
	if !strings.Contains(q, "event_type=file.uploaded") {
		t.Errorf("expected event_type filter in query, got %q", q)
	}
	if !strings.Contains(q, "from=2026-01-01") {
		t.Errorf("expected from filter in query, got %q", q)
	}
	if !strings.Contains(q, "to=2026-03-01") {
		t.Errorf("expected to filter in query, got %q", q)
	}
	if !strings.Contains(q, "page=1") {
		t.Errorf("expected page in query, got %q", q)
	}
	if !strings.Contains(q, "limit=50") {
		t.Errorf("expected limit in query, got %q", q)
	}
}

// ---------------------------------------------------------------------------
// Keys service
// ---------------------------------------------------------------------------

func TestKeysList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /keys": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[APIKey]{
				Data: []APIKey{
					{ID: "key-1", Name: "Production Key", Prefix: "cvk_live", Environment: "live", Scopes: []string{"files:read"}},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	keys, err := c.Keys.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	if keys[0].Name != "Production Key" {
		t.Errorf("expected key name 'Production Key', got %q", keys[0].Name)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/keys" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /keys": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[CreateAPIKeyResponse]{
				Data: CreateAPIKeyResponse{
					APIKey: APIKey{ID: "key-new", Name: "New Key", Prefix: "cvk_test", Environment: "test", Scopes: []string{"files:read", "files:write"}},
					Key:    "cvk_test_full_secret_key_value",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	resp, err := c.Keys.Create(context.Background(), CreateAPIKeyRequest{
		Name:        "New Key",
		Environment: "test",
		Scopes:      []string{"files:read", "files:write"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "key-new" {
		t.Errorf("expected key ID 'key-new', got %q", resp.ID)
	}
	if resp.Key != "cvk_test_full_secret_key_value" {
		t.Errorf("expected full key value, got %q", resp.Key)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/keys" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysGet(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /keys/key-1": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[APIKey]{
				Data: APIKey{ID: "key-1", Name: "My Key", Prefix: "cvk_live"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	key, err := c.Keys.Get(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.ID != "key-1" {
		t.Errorf("expected key ID 'key-1', got %q", key.ID)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/keys/key-1" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysRevoke(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /keys/key-dead": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Keys.Revoke(context.Background(), "key-dead")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/keys/key-dead" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysRotate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /keys/key-1/rotate": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[CreateAPIKeyResponse]{
				Data: CreateAPIKeyResponse{
					APIKey: APIKey{ID: "key-1", Name: "Rotated Key"},
					Key:    "cvk_live_new_rotated_value",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	resp, err := c.Keys.Rotate(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Key != "cvk_live_new_rotated_value" {
		t.Errorf("expected rotated key value, got %q", resp.Key)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/keys/key-1/rotate" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysInstantRevoke(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /api-keys/key-urgent/revoke": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Keys.InstantRevoke(context.Background(), "key-urgent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/api-keys/key-urgent/revoke" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestKeysGetRevocationStatus(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /api-keys/key-1/revocation-status": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[KeyRevocationStatus]{
				Data: KeyRevocationStatus{
					KeyID:     "key-1",
					Revoked:   true,
					RevokedAt: "2026-03-04T12:00:00Z",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	status, err := c.Keys.GetRevocationStatus(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Revoked {
		t.Error("expected key to be revoked")
	}
	if status.RevokedAt != "2026-03-04T12:00:00Z" {
		t.Errorf("expected revoked_at timestamp, got %q", status.RevokedAt)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/api-keys/key-1/revocation-status" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Comments service
// ---------------------------------------------------------------------------

func TestCommentsCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /comments": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[Comment]{
				Data: Comment{ID: "cmt-1", FileID: "file-1", Content: "Great document!", AuthorID: "user-1"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	comment, err := c.Comments.Create(context.Background(), CreateCommentRequest{
		FileID:  "file-1",
		Content: "Great document!",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != "cmt-1" {
		t.Errorf("expected comment ID 'cmt-1', got %q", comment.ID)
	}
	if comment.Content != "Great document!" {
		t.Errorf("expected content 'Great document!', got %q", comment.Content)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/comments" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestCommentsList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /comments": func(w http.ResponseWriter, r *http.Request) {
			// Verify the file_id query parameter is present
			if r.URL.Query().Get("file_id") != "file-abc" {
				writeJSON(w, 400, map[string]string{"error": "missing file_id"})
				return
			}
			writeJSON(w, 200, ListResponse[Comment]{
				Data: []Comment{
					{ID: "cmt-1", FileID: "file-abc", Content: "First comment"},
					{ID: "cmt-2", FileID: "file-abc", Content: "Second comment"},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	comments, err := c.Comments.List(context.Background(), "file-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
	if (*logs)[0].Method != "GET" {
		t.Errorf("expected GET, got %s", (*logs)[0].Method)
	}
	if !strings.Contains((*logs)[0].Query, "file_id=file-abc") {
		t.Errorf("expected file_id query param, got %q", (*logs)[0].Query)
	}
}

func TestCommentsUpdate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"PUT /comments/cmt-1": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[Comment]{
				Data: Comment{ID: "cmt-1", FileID: "file-1", Content: "Updated content"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	comment, err := c.Comments.Update(context.Background(), "cmt-1", UpdateCommentRequest{
		Content: "Updated content",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "Updated content" {
		t.Errorf("expected content 'Updated content', got %q", comment.Content)
	}
	if (*logs)[0].Method != "PUT" || (*logs)[0].Path != "/comments/cmt-1" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestCommentsDelete(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /comments/cmt-del": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Comments.Delete(context.Background(), "cmt-del")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/comments/cmt-del" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Jobs service
// ---------------------------------------------------------------------------

func TestJobsCreate(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /jobs": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[Job]{
				Data: Job{ID: "job-1", Type: "bulk_download", Status: "pending", Progress: 0},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	job, err := c.Jobs.Create(context.Background(), CreateJobRequest{
		Type:   "bulk_download",
		Params: map[string]any{"file_ids": []string{"f1", "f2"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != "job-1" {
		t.Errorf("expected job ID 'job-1', got %q", job.ID)
	}
	if job.Type != "bulk_download" {
		t.Errorf("expected job type 'bulk_download', got %q", job.Type)
	}
	if (*logs)[0].Method != "POST" || (*logs)[0].Path != "/jobs" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestJobsList(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /jobs": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[Job]{
				Data: []Job{
					{ID: "job-1", Type: "bulk_download", Status: "completed", Progress: 100},
					{ID: "job-2", Type: "export", Status: "running", Progress: 50},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	jobs, err := c.Jobs.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}
	if jobs[1].Progress != 50 {
		t.Errorf("expected second job progress 50, got %d", jobs[1].Progress)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/jobs" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestJobsCancel(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"DELETE /jobs/job-cancel": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	err := c.Jobs.Cancel(context.Background(), "job-cancel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*logs)[0].Method != "DELETE" || (*logs)[0].Path != "/jobs/job-cancel" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Bandwidth service
// ---------------------------------------------------------------------------

func TestBandwidthGetSummary(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /bandwidth": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[BandwidthSummary]{
				Data: BandwidthSummary{
					TotalUploadBytes:   1073741824,
					TotalDownloadBytes: 2147483648,
					Period:             "2026-03",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	summary, err := c.Bandwidth.GetSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.TotalUploadBytes != 1073741824 {
		t.Errorf("expected upload bytes 1073741824, got %d", summary.TotalUploadBytes)
	}
	if summary.TotalDownloadBytes != 2147483648 {
		t.Errorf("expected download bytes 2147483648, got %d", summary.TotalDownloadBytes)
	}
	if summary.Period != "2026-03" {
		t.Errorf("expected period '2026-03', got %q", summary.Period)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/bandwidth" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

func TestBandwidthGetDaily(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /bandwidth/daily": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[DailyBandwidthStats]{
				Data: []DailyBandwidthStats{
					{Date: "2026-03-01", UploadBytes: 100000, DownloadBytes: 200000},
					{Date: "2026-03-02", UploadBytes: 150000, DownloadBytes: 300000},
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	daily, err := c.Bandwidth.GetDaily(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(daily) != 2 {
		t.Fatalf("expected 2 daily entries, got %d", len(daily))
	}
	if daily[0].Date != "2026-03-01" {
		t.Errorf("expected first date '2026-03-01', got %q", daily[0].Date)
	}
	if daily[1].UploadBytes != 150000 {
		t.Errorf("expected second upload 150000, got %d", daily[1].UploadBytes)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/bandwidth/daily" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// DataExport service
// ---------------------------------------------------------------------------

func TestDataExportExport(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /users/user-abc/export": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, DataResponse[UserDataExport]{
				Data: UserDataExport{
					UserID:    "user-abc",
					Status:    "ready",
					URL:       "https://vault.test/exports/user-abc.zip",
					ExpiresAt: "2026-03-11T00:00:00Z",
					CreatedAt: "2026-03-04T12:00:00Z",
				},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	export, err := c.DataExport.Export(context.Background(), "user-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if export.UserID != "user-abc" {
		t.Errorf("expected user ID 'user-abc', got %q", export.UserID)
	}
	if export.Status != "ready" {
		t.Errorf("expected status 'ready', got %q", export.Status)
	}
	if export.URL != "https://vault.test/exports/user-abc.zip" {
		t.Errorf("expected export URL, got %q", export.URL)
	}
	if (*logs)[0].Method != "GET" || (*logs)[0].Path != "/users/user-abc/export" {
		t.Errorf("unexpected request: %s %s", (*logs)[0].Method, (*logs)[0].Path)
	}
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

func TestError401Unauthorized(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 401, map[string]string{"error": "invalid API key"})
		},
	})
	defer srv.Close()

	c := NewClient("bad-key", WithBaseURL(srv.URL))
	_, err := c.Files.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status code 401, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "invalid API key" {
		t.Errorf("expected message 'invalid API key', got %q", apiErr.Message)
	}
}

func TestError404NotFound(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files/nonexistent": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 404, map[string]string{"error": "file not found"})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := c.Files.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestError500InternalServerError(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 500, map[string]string{"error": "internal server error"})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := c.Files.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status code 500, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "internal server error" {
		t.Errorf("expected message 'internal server error', got %q", apiErr.Message)
	}
	// 500 should NOT match IsNotFound
	if IsNotFound(err) {
		t.Error("500 should not be considered a not found error")
	}
	// 500 should NOT match IsRateLimited
	if IsRateLimited(err) {
		t.Error("500 should not be considered a rate limit error")
	}
}

func TestRateLimiting(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "30")
			writeJSON(w, 429, map[string]string{"error": "rate limited"})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := c.Files.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimited(err) {
		t.Errorf("expected rate limit error, got: %v", err)
	}
	rlErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected *RateLimitError, got %T", err)
	}
	if rlErr.RetryAfter != 30*time.Second {
		t.Errorf("expected RetryAfter 30s, got %v", rlErr.RetryAfter)
	}
}

func TestRateLimitingDefaultRetryAfter(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			// No Retry-After header; should default to 60s
			writeJSON(w, 429, map[string]string{"error": "rate limited"})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := c.Files.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	rlErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected *RateLimitError, got %T", err)
	}
	if rlErr.RetryAfter != 60*time.Second {
		t.Errorf("expected default RetryAfter 60s, got %v", rlErr.RetryAfter)
	}
}

func TestErrorMessageString(t *testing.T) {
	apiErr := &APIError{StatusCode: 403, Message: "forbidden"}
	expected := "conformvault: HTTP 403: forbidden"
	if apiErr.Error() != expected {
		t.Errorf("expected %q, got %q", expected, apiErr.Error())
	}

	rlErr := &RateLimitError{
		APIError:   APIError{StatusCode: 429, Message: "rate limited"},
		RetryAfter: 45 * time.Second,
	}
	got := rlErr.Error()
	if !strings.Contains(got, "rate limited") {
		t.Errorf("expected rate limited in error string, got %q", got)
	}
	if !strings.Contains(got, "45s") {
		t.Errorf("expected 45s in error string, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Error response with non-JSON body
// ---------------------------------------------------------------------------

func TestErrorNonJSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	_, err := c.Files.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 502 {
		t.Errorf("expected status code 502, got %d", apiErr.StatusCode)
	}
	// When body is not valid JSON, the raw text becomes the message
	if apiErr.Message != "Bad Gateway" {
		t.Errorf("expected message 'Bad Gateway', got %q", apiErr.Message)
	}
}

// ---------------------------------------------------------------------------
// Webhook signature verification
// ---------------------------------------------------------------------------

func testComputeHMAC(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestWebhookVerify(t *testing.T) {
	secret := "whsec_test123"
	payload := []byte(`{"event":"file.uploaded","file_id":"f1"}`)

	// Valid signature
	validSig := testComputeHMAC(payload, secret)
	if !VerifyWebhookSignature(payload, validSig, secret) {
		t.Error("expected valid signature to pass verification")
	}

	// Invalid signature
	if VerifyWebhookSignature(payload, "bad_signature", secret) {
		t.Error("expected invalid signature to fail verification")
	}

	// Wrong secret
	if VerifyWebhookSignature(payload, validSig, "wrong-secret") {
		t.Error("expected wrong secret to fail verification")
	}
}

// ---------------------------------------------------------------------------
// Content-Type header is set only when body is present
// ---------------------------------------------------------------------------

func TestContentTypeOnlyWithBody(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[File]{Data: []File{}})
		},
		"POST /folders": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[Folder]{
				Data: Folder{ID: "fld-1", Name: "Test"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))

	// GET request should NOT have Content-Type
	_, _ = c.Files.List(context.Background(), nil)
	if (*logs)[0].CT != "" {
		t.Errorf("GET request should not have Content-Type, got %q", (*logs)[0].CT)
	}

	// POST request with body SHOULD have Content-Type
	_, _ = c.Folders.Create(context.Background(), CreateFolderRequest{Name: "Test"})
	if (*logs)[1].CT != "application/json" {
		t.Errorf("POST request with body should have Content-Type application/json, got %q", (*logs)[1].CT)
	}
}

// ---------------------------------------------------------------------------
// Request body is correctly serialized
// ---------------------------------------------------------------------------

func TestRequestBodySerialization(t *testing.T) {
	srv, logs := newMockServer(t, map[string]http.HandlerFunc{
		"POST /comments": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 201, DataResponse[Comment]{
				Data: Comment{ID: "cmt-1", FileID: "file-1", Content: "Hello"},
			})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, _ = c.Comments.Create(context.Background(), CreateCommentRequest{
		FileID:  "file-1",
		Content: "Hello",
	})

	body := (*logs)[0].Body
	if !strings.Contains(body, `"file_id":"file-1"`) {
		t.Errorf("expected file_id in body, got %q", body)
	}
	if !strings.Contains(body, `"content":"Hello"`) {
		t.Errorf("expected content in body, got %q", body)
	}
}

// ---------------------------------------------------------------------------
// doRaw error handling
// ---------------------------------------------------------------------------

func TestDoRawErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 404, map[string]string{"error": "file not found"})
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	_, err := c.Files.Download(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// doRaw success returns a body reader
// ---------------------------------------------------------------------------

func TestDoRawSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files/file-1/download" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(200)
		w.Write([]byte("fake-pdf-content"))
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	body, err := c.Files.Download(context.Background(), "file-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	b, _ := io.ReadAll(body)
	if string(b) != "fake-pdf-content" {
		t.Errorf("expected 'fake-pdf-content', got %q", string(b))
	}
}

// ---------------------------------------------------------------------------
// Empty list returns empty slice, not nil
// ---------------------------------------------------------------------------

func TestEmptyListResponse(t *testing.T) {
	srv, _ := newMockServer(t, map[string]http.HandlerFunc{
		"GET /files": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 200, ListResponse[File]{Data: []File{}})
		},
	})
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	files, err := c.Files.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if files == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}
