package campaign_test

import (
	"context"
	"database/sql"
	"testing"

	"targeting-engine/internal/campaign"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := campaign.NewPostgresRepository(sqlx.NewDb(db, "sqlmock"))

	tests := []struct {
		name        string
		mockExpect  func()
		id          string
		expected    *campaign.Campaign
		expectError bool
	}{
		{
			name: "Success - campaign found",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "image_url", "cta", "status"}).
					AddRow("test1", "Test Campaign", "http://test.com/img", "Install", "ACTIVE")
				mock.ExpectQuery(`SELECT \* FROM campaigns WHERE id = \$1`).
					WithArgs("test1").
					WillReturnRows(rows)
			},
			id: "test1",
			expected: &campaign.Campaign{
				ID:       "test1",
				Name:     "Test Campaign",
				ImageURL: "http://test.com/img",
				CTA:      "Install",
				Status:   campaign.StatusActive,
			},
		},
		{
			name: "Success - campaign not found",
			mockExpect: func() {
				mock.ExpectQuery(`SELECT \* FROM campaigns WHERE id = \$1`).
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			id:          "nonexistent",
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			result, err := repo.GetByID(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
