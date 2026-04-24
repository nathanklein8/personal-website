package repository

import (
	"context"
	"database/sql"

	"nklein.xyz/backend/models"
)

type PhotoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (r *PhotoRepository) GetAll(ctx context.Context) ([]models.Photo, error) {
	query := `SELECT id, title, file_path, alt_text, date_taken, location,
	          camera, lens, aperture, shutter_speed, iso, visible, featured, sort_order,
	          source_path
	          FROM photos ORDER BY sort_order ASC, id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := []models.Photo{}
	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.FilePath,
			&p.AltText,
			&p.DateTaken,
			&p.Location,
			&p.Camera,
			&p.Lens,
			&p.Aperture,
			&p.ShutterSpeed,
			&p.ISO,
			&p.Visible,
			&p.Featured,
			&p.SortOrder,
			&p.SourcePath,
		); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}

func (r *PhotoRepository) GetVisible(ctx context.Context) ([]models.Photo, error) {
	query := `SELECT id, title, file_path, alt_text, date_taken, location,
	          camera, lens, aperture, shutter_speed, iso, visible, featured, sort_order,
	          source_path
	          FROM photos WHERE visible = true ORDER BY sort_order ASC, id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := []models.Photo{}
	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.FilePath,
			&p.AltText,
			&p.DateTaken,
			&p.Location,
			&p.Camera,
			&p.Lens,
			&p.Aperture,
			&p.ShutterSpeed,
			&p.ISO,
			&p.Visible,
			&p.Featured,
			&p.SortOrder,
			&p.SourcePath,
		); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}

func (r *PhotoRepository) GetFeatured(ctx context.Context) ([]models.Photo, error) {
	query := `SELECT id, title, file_path, alt_text, date_taken, location,
	          camera, lens, aperture, shutter_speed, iso, visible, featured, sort_order,
	          source_path
	          FROM photos WHERE visible = true AND featured = true ORDER BY sort_order ASC, id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := []models.Photo{}
	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.FilePath,
			&p.AltText,
			&p.DateTaken,
			&p.Location,
			&p.Camera,
			&p.Lens,
			&p.Aperture,
			&p.ShutterSpeed,
			&p.ISO,
			&p.Visible,
			&p.Featured,
			&p.SortOrder,
			&p.SourcePath,
		); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}

func (r *PhotoRepository) GetByID(ctx context.Context, id int) (*models.Photo, error) {
	query := `SELECT id, title, file_path, alt_text, date_taken, location,
	          camera, lens, aperture, shutter_speed, iso, visible, featured, sort_order,
	          source_path
	          FROM photos WHERE id = $1`

	var p models.Photo
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.FilePath,
		&p.AltText,
		&p.DateTaken,
		&p.Location,
		&p.Camera,
		&p.Lens,
		&p.Aperture,
		&p.ShutterSpeed,
		&p.ISO,
		&p.Visible,
		&p.Featured,
		&p.SortOrder,
		&p.SourcePath,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PhotoRepository) Create(ctx context.Context, p *models.Photo) error {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO photos (title, file_path, alt_text, date_taken, location, camera, lens, aperture, shutter_speed, iso, visible, featured, sort_order, source_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`,
		p.Title,
		p.FilePath,
		p.AltText,
		p.DateTaken,
		p.Location,
		p.Camera,
		p.Lens,
		p.Aperture,
		p.ShutterSpeed,
		p.ISO,
		p.Visible,
		p.Featured,
		p.SortOrder,
		p.SourcePath,
	).Scan(&id)

	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (r *PhotoRepository) Update(ctx context.Context, p *models.Photo) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE photos
		SET title = $1, file_path = $2, alt_text = $3, date_taken = $4,
			location = $5, camera = $6, lens = $7, aperture = $8,
			shutter_speed = $9, iso = $10, visible = $11, featured = $12, sort_order = $13,
			source_path = $14
		WHERE id = $15
	`,
		p.Title,
		p.FilePath,
		p.AltText,
		p.DateTaken,
		p.Location,
		p.Camera,
		p.Lens,
		p.Aperture,
		p.ShutterSpeed,
		p.ISO,
		p.Visible,
		p.Featured,
		p.SortOrder,
		p.SourcePath,
		p.ID,
	)

	return err
}

func (r *PhotoRepository) UpdateTitleAndOrder(ctx context.Context, id int, title string, sortOrder int, visible *bool) error {
	if visible != nil {
		_, err := r.db.ExecContext(ctx, `
			UPDATE photos SET title = $1, sort_order = $2, visible = $3 WHERE id = $4
		`, title, sortOrder, *visible, id)
		return err
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE photos SET title = $1, sort_order = $2 WHERE id = $3
	`, title, sortOrder, id)
	return err
}

func (r *PhotoRepository) DeleteByID(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PhotoRepository) GetPhotoSourcePath(ctx context.Context, id int) (string, error) {
	var sourcePath string
	err := r.db.QueryRowContext(ctx, `SELECT source_path FROM photos WHERE id = $1`, id).Scan(&sourcePath)
	if err != nil {
		return "", err
	}
	return sourcePath, nil
}

func (r *PhotoRepository) GetAllSourcePaths(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT source_path FROM photos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}

	return paths, rows.Err()
}


