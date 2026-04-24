package models

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type Photo struct {
	ID           int     `json:"id,omitempty"`
	Title        string  `json:"title"`
	FilePath     string  `json:"filePath"`
	AltText      *string `json:"altText,omitempty"`
	DateTaken    *string `json:"dateTaken,omitempty"`
	Location     *string `json:"location,omitempty"`
	Camera       *string `json:"camera,omitempty"`
	Lens         *string `json:"lens,omitempty"`
	Aperture     *string `json:"aperture,omitempty"`
	ShutterSpeed *string `json:"shutterSpeed,omitempty"`
	ISO          *string `json:"iso,omitempty"`
	Visible      bool    `json:"visible"`
	Featured     bool    `json:"featured"`
	SortOrder    int     `json:"sortOrder"`
	SourcePath   string  `json:"sourcePath"`
	PreviewURL   string  `json:"previewURL"`
	ThumbURL     string  `json:"thumbURL"`
}

// ThumbnailPath returns the computed thumbnail path for this photo.
func (p *Photo) ThumbnailPath() string {
	return p.thumbMedPath("_thumb.jpg")
}

// MediumPath returns the computed medium preview path for this photo.
func (p *Photo) MediumPath() string {
	return p.thumbMedPath("_med.jpg")
}

// thumbMedPath builds the thumbnail or medium preview path by inserting
// the suffix before the file extension.
func (p *Photo) thumbMedPath(suffix string) string {
	ext := filepath.Ext(p.FilePath)
	return p.SourcePath + suffix + ext
}

// SetPreviewURLs populates the transient preview URL fields.
// Call this after loading a photo from the DB, before serializing to JSON.
func (p *Photo) SetPreviewURLs() {
	p.PreviewURL = fmt.Sprintf("/api/photos/available/%s/%s/%s/preview?size=med",
		url.QueryEscape(p.Year()),
		url.QueryEscape(p.Event()),
		url.QueryEscape(p.Filename()))
	p.ThumbURL = fmt.Sprintf("/api/photos/available/%s/%s/%s/preview?size=thumb",
		url.QueryEscape(p.Year()),
		url.QueryEscape(p.Event()),
		url.QueryEscape(p.Filename()))
}

// Year extracts the year component from SourcePath.
func (p *Photo) Year() string {
	parts := strings.SplitN(p.SourcePath, "/", 3)
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// Event extracts the event component from SourcePath.
func (p *Photo) Event() string {
	parts := strings.SplitN(p.SourcePath, "/", 3)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// Filename extracts the filename component from SourcePath.
func (p *Photo) Filename() string {
	parts := strings.SplitN(p.SourcePath, "/", 3)
	if len(parts) >= 3 {
		return parts[2]
	}
	return p.FilePath
}

type PhotoCreateRequest struct {
	Title     string `json:"title"`
	Filename  string `json:"filename"`
	SortOrder int    `json:"sortOrder"`
}

type PhotoUpdateRequest struct {
	ID        int     `json:"id,omitempty"`
	Title     *string `json:"title,omitempty"`
	SortOrder *int    `json:"sortOrder,omitempty"`
	Visible   *bool   `json:"visible,omitempty"`
	Featured  *bool   `json:"featured,omitempty"`
}

// DirectoryEntry represents a directory listing for the file explorer
type DirectoryEntry struct {
	Name string `json:"name"`
}

// AddPhotoRequest is used for adding a photo from the file explorer
type AddPhotoRequest struct {
	Filename  string `json:"filename"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
}
