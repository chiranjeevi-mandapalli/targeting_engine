package targeting

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RuleRepository interface {
	GetByCampaignID(ctx context.Context, campaignID string) ([]Rule, error)
	GetByCampaignIDs(ctx context.Context, campaignIDs []string) ([]Rule, error)
	Store(ctx context.Context, rule *Rule) error
	DeleteByCampaign(ctx context.Context, campaignID string) error
}

type PostgresRuleRepository struct {
	db *sqlx.DB
}

func NewPostgresRuleRepository(db *sqlx.DB) *PostgresRuleRepository {
	return &PostgresRuleRepository{db: db}
}

func (r *PostgresRuleRepository) GetByCampaignID(ctx context.Context, campaignID string) ([]Rule, error) {
	const query = `SELECT * FROM targeting_rules WHERE campaign_id = $1`
	var rules []Rule
	err := r.db.SelectContext(ctx, &rules, query, campaignID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting rules: %w", err)
	}
	return rules, nil
}

func (r *PostgresRuleRepository) GetByCampaignIDs(ctx context.Context, campaignIDs []string) ([]Rule, error) {
	if len(campaignIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT campaign_id, dimension, operation, values FROM targeting_rules WHERE campaign_id IN (?)`, campaignIDs)
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	query = r.db.Rebind(query)

	var rules []Rule
	err = r.db.SelectContext(ctx, &rules, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting rules: %w", err)
	}

	return rules, nil
}

func (r *PostgresRuleRepository) Store(ctx context.Context, rule *Rule) error {
	const query = `
		INSERT INTO targeting_rules (campaign_id, dimension, operation, values)
		VALUES (:campaign_id, :dimension, :operation, :values)
	`
	_, err := r.db.NamedExecContext(ctx, query, rule)
	if err != nil {
		return fmt.Errorf("error storing rule: %w", err)
	}
	return nil
}

func (r *PostgresRuleRepository) DeleteByCampaign(ctx context.Context, campaignID string) error {
	const query = `DELETE FROM targeting_rules WHERE campaign_id = $1`
	_, err := r.db.ExecContext(ctx, query, campaignID)
	if err != nil {
		return fmt.Errorf("error deleting rules: %w", err)
	}
	return nil
}
