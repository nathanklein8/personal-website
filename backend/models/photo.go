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
}

type PhotoCreateRequest struct {
	Title        *string   `json:"title"`
	FilePath     string    `json:"file_path"`
	AltText      *string   `json:"altText,omitempty"`
	DateTaken    *string   `json:"dateTaken,omitempty"`
	Location     *string   `json:"location,omitempty"`
	Camera       *string   `json:"camera,omitempty"`
	Lens         *string   `json:"lens,omitempty"`
	Aperture     *string   `json:"aperture,omitempty"`
	ShutterSpeed *string   `json:"shutterSpeed,omitempty"`
	ISO          *string   `json:"iso,omitempty"`
	Visible      *bool     `json:"visible,omitempty"`
	Featured     *bool     `json:"featured,omitempty"`
	SortOrder    *int      `json:"sortOrder,omitempty"`
}

type PhotoUpdateRequest struct {
	ID           int       `json:"id,omitempty"`
	Title        *string   `json:"title"`
	FilePath     string    `json:"file_path,omitempty"`
	AltText      *string   `json:"altText,omitempty"`
	DateTaken    *string   `json:"dateTaken,omitempty"`
	Location     *string   `json:"location,omitempty"`
	Camera       *string   `json:"camera,omitempty"`
	Lens         *string   `json:"lens,omitempty"`
	Aperture     *string   `json:"aperture,omitempty"`
	ShutterSpeed *string   `json:"shutterSpeed,omitempty"`
	ISO          *string   `json:"iso,omitempty"`
	Visible      bool      `json:"visible,omitempty"`
	Featured     bool      `json:"featured,omitempty"`
	SortOrder    int       `json:"sortOrder,omitempty"`
}
