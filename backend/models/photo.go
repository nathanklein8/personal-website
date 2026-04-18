package models

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
	Featured     bool    `json:"featured,omitempty"`
	SortOrder    int     `json:"sortOrder"`
	SourcePath   string  `json:"sourcePath"`
	ThumbnailPath *string `json:"thumbnailPath,omitempty"`
	MediumPath   *string `json:"mediumPath,omitempty"`
}

type PhotoCreateRequest struct {
	Title       string `json:"title"`
	Filename    string `json:"filename"`
	SortOrder   int    `json:"sortOrder"`
}

type PhotoUpdateRequest struct {
	ID        int    `json:"id,omitempty"`
	Title     *string `json:"title,omitempty"`
	SortOrder *int    `json:"sortOrder,omitempty"`
	Visible   *bool   `json:"visible,omitempty"`
}

// DirectoryEntry represents a directory listing for the file explorer
type DirectoryEntry struct {
	Name string `json:"name"`
}

// AddPhotoRequest is used for adding a photo from the file explorer
type AddPhotoRequest struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
	SortOrder int   `json:"sortOrder"`
}
