package seed

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"targeting-engine/internal/campaign"
	"targeting-engine/internal/targeting"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	if err := s.SeedCampaigns(ctx); err != nil {
		return fmt.Errorf("seeding campaigns: %w", err)
	}
	if err := s.SeedTargetingRules(ctx); err != nil {
		return fmt.Errorf("seeding targeting rules: %w", err)
	}
	return nil
}

func (s *Seeder) SeedCampaigns(ctx context.Context) error {
	campaigns := []campaign.Campaign{
		{
			ID:       "spotify",
			Name:     "Spotify - Music for everyone",
			ImageURL: "https://somelink",
			CTA:      "Download",
			Status:   campaign.StatusActive,
		},
		{
			ID:       "duolingo",
			Name:     "Duolingo: Best way to learn",
			ImageURL: "https://somelink2",
			CTA:      "Install",
			Status:   campaign.StatusActive,
		},
		{
			ID:       "subwaysurfer",
			Name:     "Subway Surfer",
			ImageURL: "https://somelink3",
			CTA:      "Play",
			Status:   campaign.StatusActive,
		},
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, "DELETE FROM campaigns"); err != nil {
		return fmt.Errorf("clearing campaigns: %w", err)
	}

	for _, c := range campaigns {
		query := `INSERT INTO campaigns (id, name, image_url, cta, status, created_at, updated_at) 
		          VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = tx.ExecContext(ctx, query,
			c.ID,
			c.Name,
			c.ImageURL,
			c.CTA,
			string(c.Status),
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("inserting campaign %s: %w", c.ID, err)
		}
	}

	return tx.Commit()
}

func (s *Seeder) SeedTargetingRules(ctx context.Context) error {
	raw := func(values []string) json.RawMessage {
		bytes, _ := json.Marshal(values)
		return json.RawMessage(bytes)
	}

	rules := []targeting.Rule{
		{
			CampaignID: "spotify",
			Dimension:  targeting.DimensionCountry,
			Operation:  targeting.OperationInclude,
			Values:     raw([]string{"US", "Canada"}),
		},
		{
			CampaignID: "duolingo",
			Dimension:  targeting.DimensionOS,
			Operation:  targeting.OperationInclude,
			Values:     raw([]string{"Android", "iOS"}),
		},
		{
			CampaignID: "duolingo",
			Dimension:  targeting.DimensionCountry,
			Operation:  targeting.OperationExclude,
			Values:     raw([]string{"US"}),
		},
		{
			CampaignID: "subwaysurfer",
			Dimension:  targeting.DimensionOS,
			Operation:  targeting.OperationInclude,
			Values:     raw([]string{"Android"}),
		},
		{
			CampaignID: "subwaysurfer",
			Dimension:  targeting.DimensionApp,
			Operation:  targeting.OperationInclude,
			Values:     raw([]string{"com.gametion.ludokinggame"}),
		},
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, "DELETE FROM targeting_rules"); err != nil {
		return fmt.Errorf("clearing targeting rules: %w", err)
	}

	for _, r := range rules {
		query := `INSERT INTO targeting_rules (campaign_id, dimension, operation, values) 
		          VALUES ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, query,
			r.CampaignID,
			string(r.Dimension),
			string(r.Operation),
			r.Values, // Already json.RawMessage (a []byte)
		)
		if err != nil {
			return fmt.Errorf("inserting rule for campaign %s: %w", r.CampaignID, err)
		}
	}

	return tx.Commit()
}
