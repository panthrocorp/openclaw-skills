package google

import (
	"context"
	"fmt"
	"io"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// googleAppsMIMEExports maps Google Workspace MIME types to their plain-text
// export equivalents. Files with these MIME types cannot be downloaded directly
// and must be exported via Files.Export instead.
var googleAppsMIMEExports = map[string]string{
	"application/vnd.google-apps.document":     "text/plain",
	"application/vnd.google-apps.spreadsheet":  "text/csv",
	"application/vnd.google-apps.presentation": "text/plain",
}

// DriveClient provides read-only access to the Google Drive API.
// No create, update, or delete operations exist in this type.
type DriveClient struct {
	svc *drive.Service
}

// NewDriveClient creates a Drive API client using the provided token source.
func NewDriveClient(ctx context.Context, ts oauth2.TokenSource) (*DriveClient, error) {
	svc, err := drive.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("creating drive service: %w", err)
	}
	return &DriveClient{svc: svc}, nil
}

// driveFileFields is the default field mask for file list responses.
const driveFileFields = "files(id,name,mimeType,size,createdTime,modifiedTime,parents,webViewLink)"

// ListFiles returns files visible to the authenticated user.
func (c *DriveClient) ListFiles(ctx context.Context, query string, pageSize int64) ([]*drive.File, error) {
	call := c.svc.Files.List().
		PageSize(pageSize).
		Fields(driveFileFields).
		OrderBy("modifiedTime desc").
		Context(ctx)

	if query != "" {
		call = call.Q(query)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("listing files: %w", err)
	}
	return resp.Files, nil
}

// GetFile retrieves metadata for a single file.
func (c *DriveClient) GetFile(ctx context.Context, fileID string) (*drive.File, error) {
	file, err := c.svc.Files.Get(fileID).
		Fields("*").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("getting file %s: %w", fileID, err)
	}
	return file, nil
}

// DownloadFile retrieves the content of a file. For Google Workspace files
// (Docs, Sheets, Slides) it exports to a text format automatically. For all
// other files it downloads the raw bytes. The caller must close the returned
// ReadCloser.
func (c *DriveClient) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	meta, err := c.svc.Files.Get(fileID).
		Fields("id,mimeType").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("getting file metadata %s: %w", fileID, err)
	}

	if exportMIME, ok := googleAppsMIMEExports[meta.MimeType]; ok {
		resp, err := c.svc.Files.Export(fileID, exportMIME).Context(ctx).Download()
		if err != nil {
			return nil, fmt.Errorf("exporting file %s as %s: %w", fileID, exportMIME, err)
		}
		return resp.Body, nil
	}

	resp, err := c.svc.Files.Get(fileID).Context(ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("downloading file %s: %w", fileID, err)
	}
	return resp.Body, nil
}

// IsGoogleAppsFile reports whether the MIME type is a Google Workspace type
// that requires export rather than direct download.
func IsGoogleAppsFile(mimeType string) bool {
	return strings.HasPrefix(mimeType, "application/vnd.google-apps.")
}
