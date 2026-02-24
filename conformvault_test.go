package conformvault

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func TestFilesList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("User-Agent") != userAgent {
			t.Errorf("unexpected User-Agent: %s", r.Header.Get("User-Agent"))
		}

		resp := ListResponse[File]{
			Data: []File{
				{ID: "f1", Name: "test.pdf", Size: 1024},
				{ID: "f2", Name: "doc.txt", Size: 512},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
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
}

func TestRateLimiting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "rate limited"})
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
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

func TestNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "file not found"})
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	_, err := c.Files.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestSignatureCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/signatures" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req CreateSignatureRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Subject != "NDA Agreement" {
			t.Errorf("expected subject 'NDA Agreement', got %q", req.Subject)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(DataResponse[SignatureEnvelope]{
			Data: SignatureEnvelope{
				ID:      "env-123",
				Status:  "sent",
				Subject: req.Subject,
			},
		})
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
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
}

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
