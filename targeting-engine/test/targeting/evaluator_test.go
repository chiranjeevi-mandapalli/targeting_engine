package targeting_test

import (
	"context"
	"testing"

	"targeting-engine/internal/targeting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRuleRepository implements targeting.RuleRepository
type MockRuleRepository struct {
	mock.Mock
}

func (m *MockRuleRepository) GetByCampaignID(ctx context.Context, campaignID string) ([]targeting.Rule, error) {
	args := m.Called(ctx, campaignID)
	return args.Get(0).([]targeting.Rule), args.Error(1)
}

func (m *MockRuleRepository) GetByCampaignIDs(ctx context.Context, campaignIDs []string) ([]targeting.Rule, error) {
	args := m.Called(ctx, campaignIDs)
	return args.Get(0).([]targeting.Rule), args.Error(1)
}

func (m *MockRuleRepository) Store(ctx context.Context, rule *targeting.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRuleRepository) DeleteByCampaign(ctx context.Context, campaignID string) error {
	args := m.Called(ctx, campaignID)
	return args.Error(0)
}

func TestEvaluator_Evaluate(t *testing.T) {
	mockRepo := new(MockRuleRepository)
	evaluator := targeting.NewEvaluator(mockRepo)

	tests := []struct {
		name           string
		mockRules      []targeting.Rule
		app            string
		country        string
		os             string
		campaignIDs    []string
		expectedResult []string
	}{
		{
			name: "Include country match",
			mockRules: []targeting.Rule{
				{
					CampaignID: "camp1",
					Dimension:  targeting.DimensionCountry,
					Operation:  targeting.OperationInclude,
					Values:     []string{"US"},
				},
			},
			country:        "US",
			campaignIDs:    []string{"camp1"},
			expectedResult: []string{"camp1"},
		},
		{
			name: "Exclude country match",
			mockRules: []targeting.Rule{
				{
					CampaignID: "camp1",
					Dimension:  targeting.DimensionCountry,
					Operation:  targeting.OperationExclude,
					Values:     []string{"US"},
				},
			},
			country:        "US",
			campaignIDs:    []string{"camp1"},
			expectedResult: []string{},
		},
		{
			name: "Multiple rules with one match",
			mockRules: []targeting.Rule{
				{
					CampaignID: "camp1",
					Dimension:  targeting.DimensionCountry,
					Operation:  targeting.OperationInclude,
					Values:     []string{"US"},
				},
				{
					CampaignID: "camp1",
					Dimension:  targeting.DimensionOS,
					Operation:  targeting.OperationExclude,
					Values:     []string{"iOS"},
				},
			},
			country:        "US",
			os:             "Android",
			campaignIDs:    []string{"camp1"},
			expectedResult: []string{"camp1"},
		},
		{
			name:           "No rules for campaign",
			mockRules:      []targeting.Rule{},
			country:        "US",
			campaignIDs:    []string{"camp1"},
			expectedResult: []string{"camp1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			mockRepo.On("GetByCampaignIDs", mock.Anything, tt.campaignIDs).Return(tt.mockRules, nil)

			// Execute evaluation
			result, err := evaluator.Evaluate(context.Background(), tt.app, tt.country, tt.os, tt.campaignIDs)

			// Verify results
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.expectedResult, result)

			// Verify mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
