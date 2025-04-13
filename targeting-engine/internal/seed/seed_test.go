package seed_test

import (
	"context"
	"database/sql"
	"testing"

	"targeting-engine/internal/seed"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSeeder(t *testing.T) {
	// Setup test database connection
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create clean test environment
	_, err = db.Exec("DROP TABLE IF EXISTS campaigns, targeting_rules")
	if err != nil {
		t.Fatalf("Failed to clean test database: %v", err)
	}

	// Create tables (you might want to run your migrations here)
	_, err = db.Exec(`
		CREATE TABLE campaigns (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			image_url TEXT NOT NULL,
			cta VARCHAR(255) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		);

		CREATE TABLE targeting_rules (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			campaign_id VARCHAR(255) NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
			dimension VARCHAR(50) NOT NULL,
			operation VARCHAR(50) NOT NULL,
			values JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Test seeding
	t.Run("SeedAll", func(t *testing.T) {
		seeder := seed.NewSeeder(db)
		ctx := context.Background()

		err := seeder.SeedAll(ctx)
		assert.NoError(t, err)

		// Verify campaigns were seeded
		var campaignCount int
		err = db.QueryRow("SELECT COUNT(*) FROM campaigns").Scan(&campaignCount)
		assert.NoError(t, err)
		assert.Equal(t, 3, campaignCount)

		// Verify rules were seeded
		var ruleCount int
		err = db.QueryRow("SELECT COUNT(*) FROM targeting_rules").Scan(&ruleCount)
		assert.NoError(t, err)
		assert.Equal(t, 5, ruleCount)
	})
}
