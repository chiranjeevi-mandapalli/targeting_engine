package campaign

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Campaign, error)
	GetActive(ctx context.Context) ([]Campaign, error)
	GetByIDs(ctx context.Context, ids []string) ([]Campaign, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Campaign, error) {
	const query = `SELECT * FROM campaigns WHERE id = $1`
	var c Campaign
	err := r.db.GetContext(ctx, &c, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting campaign: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepository) GetActive(ctx context.Context) ([]Campaign, error) {
	const query = `SELECT * FROM campaigns WHERE status = 'ACTIVE'`
	var campaigns []Campaign
	err := r.db.SelectContext(ctx, &campaigns, query)
	if err != nil {
		return nil, fmt.Errorf("error getting active campaigns: %w", err)
	}
	return campaigns, nil
}

func (r *PostgresRepository) GetByIDs(ctx context.Context, ids []string) ([]Campaign, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM campaigns WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	query = r.db.Rebind(query)
	var campaigns []Campaign
	err = r.db.SelectContext(ctx, &campaigns, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting campaigns by IDs: %w", err)
	}
	return campaigns, nil
}
