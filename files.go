package conformvault

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// FilesService handles file operations.
type FilesService struct {
	client *Client
}

// List returns all files, optionally filtered.
func (s *FilesService) List(ctx context.Context, opts *FileListOptions) ([]File, error) {
	path := "/files"
	if opts != nil {
		sep := "?"
		if opts.FolderID != nil {
			path += sep + "folder_id=" + url.QueryEscape(*opts.FolderID)
			sep = "&"
		}
		if opts.Page > 0 {
			path += sep + fmt.Sprintf("page=%d", opts.Page)
			sep = "&"
		}
		if opts.Limit > 0 {
			path += sep + fmt.Sprintf("limit=%d", opts.Limit)
		}
	}

	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[File]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get returns a single file by ID.
func (s *FilesService) Get(ctx context.Context, fileID string) (*File, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID, nil)
	if err != nil {
		return nil, err
	}

	var resp DataResponse[File]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Upload uploads a file. The reader should contain the file content.
func (s *FilesService) Upload(ctx context.Context, reader io.Reader, filename string, folderID *string) (*UploadResult, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		if folderID != nil {
			writer.WriteField("folder_id", *folderID)
		}

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, reader); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	url := s.client.baseURL + "/files"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, pr)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	var resp DataResponse[UploadResult]
	if err := s.client.do(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Download returns a reader for the file content.
func (s *FilesService) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/download", nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}

// Delete deletes a file by ID.
func (s *FilesService) Delete(ctx context.Context, fileID string) error {
	req, err := s.client.newRequest(ctx, "DELETE", "/files/"+fileID, nil)
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// GetThumbnail returns the thumbnail image for a file as a raw stream.
func (s *FilesService) GetThumbnail(ctx context.Context, fileID string) (io.ReadCloser, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/thumbnail", nil)
	if err != nil {
		return nil, err
	}
	return s.client.doRaw(req)
}

// GetScanReport returns the antivirus scan report for a file.
func (s *FilesService) GetScanReport(ctx context.Context, fileID string) (*FileScanReport, error) {
	req, err := s.client.newRequest(ctx, "GET", "/files/"+fileID+"/scan-report", nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[FileScanReport]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
