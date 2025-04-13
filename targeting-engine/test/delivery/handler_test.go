package delivery_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"targeting-engine/internal/delivery"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetMatchingCampaigns(ctx context.Context, req delivery.Request) (delivery.Response, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(delivery.Response), args.Error(1)
}

func TestDeliveryHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mockService)
		expectedStatus int
	}{
		{
			name: "Valid request with matching campaigns",
			url:  "/v1/delivery?app=com.test&country=US&os=Android",
			mockSetup: func(ms *mockService) {
				ms.On("GetMatchingCampaigns", mock.Anything, delivery.Request{
					App:     "com.test",
					Country: "US",
					OS:      "Android",
				}).Return(delivery.Response{
					Campaigns: []delivery.CampaignResponse{
						{ID: "test1", Image: "http://test.com/img", CTA: "Install"},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing app parameter",
			url:            "/v1/delivery?country=US&os=Android",
			mockSetup:      func(ms *mockService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No matching campaigns",
			url:  "/v1/delivery?app=com.test&country=DE&os=iOS",
			mockSetup: func(ms *mockService) {
				ms.On("GetMatchingCampaigns", mock.Anything, delivery.Request{
					App:     "com.test",
					Country: "DE",
					OS:      "iOS",
				}).Return(delivery.Response{
					Error: delivery.ErrNoCampaigns.Error(),
				}, delivery.ErrNoCampaigns)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			tt.mockSetup(mockSvc)

			handler := delivery.MakeHTTPHandler(mockSvc)
			req := httptest.NewRequest("GET", tt.url, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}
