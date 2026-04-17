package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"nklein.xyz/backend/models"
)

type LandingCardRepository struct {
	db *sql.DB
}

func NewLandingCardRepository(db *sql.DB) *LandingCardRepository {
	return &LandingCardRepository{db: db}
}

func (r *LandingCardRepository) GetByID(ctx context.Context, id int) (*models.LandingCard, error) {
	query := `SELECT bio, email, linkedin, github, skills FROM landing_card WHERE id = $1`

	var lc models.LandingCard
	var skillsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(&lc.Bio, &lc.Email, &lc.Linkedin, &lc.Github, &skillsJSON)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(skillsJSON, &lc.Skills); err != nil {
		return nil, err
	}

	return &lc, nil
}

func (r *LandingCardRepository) CreateOrUpdate(ctx context.Context, lc *models.LandingCard) error {
	skillsJSON, err := json.Marshal(lc.Skills)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO landing_card (id, bio, email, linkedin, github, skills)
		VALUES (1, $1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			bio = EXCLUDED.bio,
			email = EXCLUDED.email,
			linkedin = EXCLUDED.linkedin,
			github = EXCLUDED.github,
			skills = EXCLUDED.skills
	`,
		lc.Bio,
		lc.Email,
		lc.Linkedin,
		lc.Github,
		skillsJSON,
	)

	return err
}
