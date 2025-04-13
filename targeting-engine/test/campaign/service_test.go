package campaign_test

import (
	"context"
	"testing"

	"targeting-engine/internal/campaign"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*campaign.Campaign, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) GetActive(ctx context.Context) ([]campaign.Campaign, error) {
	args := m.Called(ctx)
	return args.Get(0).([]campaign.Campaign), args.Error(1)
}

func (m *MockRepository) GetByIDs(ctx context.Context, ids []string) ([]campaign.Campaign, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]campaign.Campaign), args.Error(1)
}

func TestService_GetCampaignByID(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := campaign.NewService(mockRepo)

	tests := []struct {
		name        string
		setupMock   func()
		id          string
		expected    *campaign.Campaign
		expectError bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, "valid").
					Return(&campaign.Campaign{ID: "valid"}, nil)
			},
			id:       "valid",
			expected: &campaign.Campaign{ID: "valid"},
		},
		{
			name:        "Empty ID",
			setupMock:   func() {},
			id:          "",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := svc.GetCampaignByID(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
