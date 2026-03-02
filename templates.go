package conformvault

import (
	"context"
	"fmt"
	"io"
)

// TemplatesService handles document template operations.
type TemplatesService struct {
	client *Client
}

// --- Types ---

// DocumentTemplate represents a document template.
type DocumentTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	ContentType string   `json:"content_type"`
	Fields      []string `json:"fields,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// GeneratedDocument represents a document generated from a template.
type GeneratedDocument struct {
	ID         string `json:"id"`
	TemplateID string `json:"template_id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Status     string `json:"status"`
	FileID     string `json:"file_id,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// TemplateListResponse wraps a paginated list of templates.
type TemplateListResponse struct {
	Data  []DocumentTemplate `json:"data"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Total int                `json:"total"`
}

// CreateTemplateRequest is the input for creating a template.
type CreateTemplateRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	ContentType string   `json:"content_type"`
	Fields      []string `json:"fields,omitempty"`
}

// UpdateTemplateRequest is the input for updating a template.
type UpdateTemplateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Fields      []string `json:"fields,omitempty"`
}

// GenerateDocumentRequest is the input for generating a document from a template.
type GenerateDocumentRequest struct {
	Data     map[string]string `json:"data"`
	Filename *string           `json:"filename,omitempty"`
}

// --- Methods ---

// Create creates a new document template.
func (s *TemplatesService) Create(ctx context.Context, r CreateTemplateRequest) (*DocumentTemplate, error) {
	req, err := s.client.newRequest(ctx, "POST", "/templates", r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[DocumentTemplate]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of document templates.
func (s *TemplatesService) List(ctx context.Context, page, limit int) (*TemplateListResponse, error) {
	path := "/templates"
	sep := "?"
	if page > 0 {
		path += fmt.Sprintf("%spage=%d", sep, page)
		sep = "&"
	}
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp TemplateListResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single document template by ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*DocumentTemplate, error) {
	req, err := s.client.newRequest(ctx, "GET", "/templates/"+id, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[DocumentTemplate]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates a document template.
func (s *TemplatesService) Update(ctx context.Context, id string, r UpdateTemplateRequest) (*DocumentTemplate, error) {
	req, err := s.client.newRequest(ctx, "PUT", "/templates/"+id, r)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[DocumentTemplate]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a document template.
func (s *TemplatesService) Delete(ctx context.Context, id string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/templates/"+id, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// Generate generates a PDF document from a template and returns the binary stream.
// The caller is responsible for closing the returned reader.
func (s *TemplatesService) Generate(ctx context.Context, id string, r GenerateDocumentRequest) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "POST", "/templates/"+id+"/generate", r)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}

// ListDocuments returns all documents generated from a template.
func (s *TemplatesService) ListDocuments(ctx context.Context, id string) ([]GeneratedDocument, error) {
	req, err := s.client.newRequest(ctx, "GET", "/templates/"+id+"/documents", nil)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[GeneratedDocument]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
