package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
)

const (
	photoLibraryPath = "/photos"
	thumbnailPath    = "/thumbnails"
	thumbWidth       = 200
	mediumWidth      = 800
)

type PhotoService struct {
	repo *repository.PhotoRepository
}

func NewPhotoService(repo *repository.PhotoRepository) *PhotoService {
	return &PhotoService{repo: repo}
}

// --- Directory Browsing ---

func (s *PhotoService) ListYears(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(photoLibraryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read photo library: %w", err)
	}

	var years []string
	for _, e := range entries {
		if e.IsDir() {
			years = append(years, e.Name())
		}
	}
	return years, nil
}

func (s *PhotoService) ListEvents(ctx context.Context, year string) ([]string, error) {
	yearPath := filepath.Join(photoLibraryPath, year)
	entries, err := os.ReadDir(yearPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read year directory %s: %w", year, err)
	}

	var events []string
	for _, e := range entries {
		if e.IsDir() {
			events = append(events, e.Name())
		}
	}
	return events, nil
}

func (s *PhotoService) ListPhotos(ctx context.Context, year, event string) ([]string, error) {
	eventPath := filepath.Join(photoLibraryPath, year, event)
	entries, err := os.ReadDir(eventPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read event directory %s/%s: %w", year, event, err)
	}

	var photos []string
	for _, e := range entries {
		if !e.IsDir() && strings.ToLower(filepath.Ext(e.Name())) == ".jpg" {
			photos = append(photos, e.Name())
		}
	}
	return photos, nil
}

// --- Image Processing ---

func (s *PhotoService) GenerateThumbnails(sourceAbsPath, thumbnailAbsPath, mediumAbsPath string) error {
	src, err := imaging.Open(sourceAbsPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Generate thumbnail (200px wide, fit)
	thumb := imaging.Resize(src, thumbWidth, 0, imaging.Lanczos)
	if err := imaging.Save(thumb, thumbnailAbsPath); err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	// Generate medium preview (800px wide, fit)
	med := imaging.Resize(src, mediumWidth, 0, imaging.Lanczos)
	if err := imaging.Save(med, mediumAbsPath); err != nil {
		return fmt.Errorf("failed to save medium preview: %w", err)
	}

	return nil
}

func (s *PhotoService) ExtractEXIF(sourceAbsPath string) (*models.Photo, error) {
	f, err := os.Open(sourceAbsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image for EXIF: %w", err)
	}
	defer f.Close()

	exifData, err := exif.Decode(f)
	if err != nil {
		// No EXIF data is not an error
		return &models.Photo{}, nil
	}

	photo := &models.Photo{}

	// Extract date taken
	if t, err := exifData.DateTime(); err == nil {
		dateTaken := t.Format("2006-01-02")
		photo.DateTaken = &dateTaken
	}

	// Extract camera make/model
	if makeModel, err := exifData.Get(exif.Model); err == nil {
		val, _ := makeModel.StringVal()
		photo.Camera = &val
	}

	// Extract lens
	if lens, err := exifData.Get(exif.LensModel); err == nil {
		val, _ := lens.StringVal()
		photo.Lens = &val
	}

	// Extract aperture (FNumber)
	if aperture, err := exifData.Get(exif.FNumber); err == nil {
		num, den, _ := aperture.Rat2(0)
		if den != 0 {
			apertureStr := fmt.Sprintf("f/%.1f", float64(num)/float64(den))
			photo.Aperture = &apertureStr
		}
	}

	// Extract exposure time (shutter speed)
	if expTime, err := exifData.Get(exif.ExposureTime); err == nil {
		num, den, _ := expTime.Rat2(0)
		if den == 0 || num == 0 {
			photo.ShutterSpeed = strPtr("1s")
		} else if num < den {
			// Fractional: e.g., 1/250
			photo.ShutterSpeed = strPtr(fmt.Sprintf("1/%ds", den/num))
		} else {
			photo.ShutterSpeed = strPtr(fmt.Sprintf("%ds", num/den))
		}
	}

	// Extract ISO
	if iso, err := exifData.Get(exif.ISOSpeedRatings); err == nil {
		val, _ := iso.Int(0)
		isoStr := strconv.Itoa(val)
		photo.ISO = &isoStr
	}

	// Extract location (GPS)
	if lat, lon, err := exifData.LatLong(); err == nil {
		location := fmt.Sprintf("%.4f, %.4f", lat, lon)
		photo.Location = &location
	}

	return photo, nil
}

// --- CRUD Operations ---

func (s *PhotoService) GetAll(ctx context.Context) ([]models.Photo, error) {
	return s.repo.GetAll(ctx)
}

func (s *PhotoService) GetByID(ctx context.Context, id int) (*models.Photo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PhotoService) AddPhoto(ctx context.Context, req models.AddPhotoRequest) (*models.Photo, error) {
	// Parse source path from filename (year/event/filename)
	parts := strings.Split(req.Filename, "/")
	if len(parts) != 3 {
		return nil, &APIError{Status: http.StatusBadRequest, Message: "filename must be in format: year/event/filename.jpg"}
	}

	year, event, filename := parts[0], parts[1], parts[2]

	// Build paths
	sourceAbsPath := filepath.Join(photoLibraryPath, year, event, filename)
	thumbRelPath := filepath.Join(year, event, filename+"_thumb.jpg")
	mediumRelPath := filepath.Join(year, event, filename+"_med.jpg")
	thumbAbsPath := filepath.Join(thumbnailPath, thumbRelPath)
	mediumAbsPath := filepath.Join(thumbnailPath, mediumRelPath)

	// Ensure thumbnail directory exists
	thumbDir := filepath.Dir(thumbAbsPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create thumbnail directory: %w", err)
	}

	// Extract EXIF data
	exifPhoto, err := s.ExtractEXIF(sourceAbsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract EXIF: %w", err)
	}

	// Generate thumbnails
	if err := s.GenerateThumbnails(sourceAbsPath, thumbAbsPath, mediumAbsPath); err != nil {
		return nil, fmt.Errorf("failed to generate thumbnails: %w", err)
	}

	// Build photo model
	photo := &models.Photo{
		Title:         req.Title,
		FilePath:      filename,
		SourcePath:    req.Filename,
		ThumbnailPath: strPtr(thumbRelPath),
		MediumPath:    strPtr(mediumRelPath),
		SortOrder:     req.SortOrder,
		Visible:       true,
		// Copy EXIF data
		Camera:       exifPhoto.Camera,
		Lens:         exifPhoto.Lens,
		Aperture:     exifPhoto.Aperture,
		ShutterSpeed: exifPhoto.ShutterSpeed,
		ISO:          exifPhoto.ISO,
		Location:     exifPhoto.Location,
		DateTaken:    exifPhoto.DateTaken,
	}

	if err := s.repo.Create(ctx, photo); err != nil {
		return nil, fmt.Errorf("failed to save photo record: %w", err)
	}

	return photo, nil
}

func (s *PhotoService) Update(ctx context.Context, id int, req models.PhotoUpdateRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.SortOrder != nil {
		existing.SortOrder = *req.SortOrder
	}
	if req.Visible != nil {
		existing.Visible = *req.Visible
	}

	return s.repo.Update(ctx, existing)
}

func (s *PhotoService) DeleteByID(ctx context.Context, id int) error {
	// Get photo info before deleting
	photo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete generated thumbnails from thumbnails volume
	if photo.ThumbnailPath != nil {
		thumbPath := filepath.Join(thumbnailPath, *photo.ThumbnailPath)
		os.Remove(thumbPath)
	}
	if photo.MediumPath != nil {
		medPath := filepath.Join(thumbnailPath, *photo.MediumPath)
		os.Remove(medPath)
	}

	// Delete DB row
	return s.repo.DeleteByID(ctx, id)
}

func (s *PhotoService) RegenerateAllThumbnails(ctx context.Context) error {
	photos, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get photos: %w", err)
	}

	for _, photo := range photos {
		parts := strings.Split(photo.SourcePath, "/")
		if len(parts) != 3 {
			continue
		}
		year, event, filename := parts[0], parts[1], parts[2]

		sourceAbsPath := filepath.Join(photoLibraryPath, year, event, filename)
		thumbRelPath := filepath.Join(year, event, filename+"_thumb.jpg")
		mediumRelPath := filepath.Join(year, event, filename+"_med.jpg")
		thumbAbsPath := filepath.Join(thumbnailPath, thumbRelPath)
		mediumAbsPath := filepath.Join(thumbnailPath, mediumRelPath)

		// Ensure thumbnail directory exists
		thumbDir := filepath.Dir(thumbAbsPath)
		if err := os.MkdirAll(thumbDir, 0755); err != nil {
			continue
		}

		// Generate thumbnails
		if err := s.GenerateThumbnails(sourceAbsPath, thumbAbsPath, mediumAbsPath); err != nil {
			continue
		}

		// Update DB with new paths
		if err := s.repo.UpdateThumbnailPaths(ctx, photo.ID, strPtr(thumbRelPath), strPtr(mediumRelPath)); err != nil {
			continue
		}
	}

	return nil
}

// Helper
func strPtr(s string) *string {
	return &s
}
