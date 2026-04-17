package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"nklein.xyz/backend/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	query := `SELECT id, icon, title, description, technologies, deployment_link, image, alt_text
	          FROM projects ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var techJSON []byte

		if err := rows.Scan(
			&p.ID,
			&p.Icon,
			&p.Title,
			&p.Description,
			&techJSON,
			&p.DeploymentLink,
			&p.Image,
			&p.AltText,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(techJSON, &p.Technologies); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func (r *ProjectRepository) CreateOrUpdate(ctx context.Context, p *models.Project) error {
	techJSON, err := json.Marshal(p.Technologies)
	if err != nil {
		return err
	}

	if p.ID == 0 {
		var result sql.Result
		var err error
		result, err = r.db.ExecContext(ctx, `
			INSERT INTO projects (icon, title, description, technologies, deployment_link, image, alt_text)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, p.Icon, p.Title, p.Description, techJSON, p.DeploymentLink, p.Image, p.AltText)
		if err != nil {
			return err
		}
		_ = result // suppress unused variable warning
	} else {
		_, err = r.db.ExecContext(ctx, `
			UPDATE projects
			SET icon = $1, title = $2, description = $3, technologies = $4,
				deployment_link = $5, image = $6, alt_text = $7
			WHERE id = $8
		`, p.Icon, p.Title, p.Description, techJSON, p.DeploymentLink, p.Image, p.AltText, p.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) DeleteByID(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, id)
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
