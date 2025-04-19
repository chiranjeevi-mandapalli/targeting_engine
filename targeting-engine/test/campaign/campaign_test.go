package campaign_test

// import (
//     "net/http"
//     "net/http/httptest"
//     "strings"
//     "testing"

//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/mock"
//     "targeting-engine/internal/handler"
// )

// type MockCampaignService struct {
//     mock.Mock
// }

// func (m *MockCampaignService) CreateCampaign(data map[string]interface{}) error {
//     args := m.Called(data)
//     return args.Error(0)
// }

// func (m *MockCampaignService) GetCampaigns() ([]map[string]interface{}, error) {
//     args := m.Called()
//     return args.Get(0).([]map[string]interface{}), args.Error(1)
// }

// func TestCreateCampaignMock(t *testing.T) {
//     mockService := new(MockCampaignService)
//     mockService.On("CreateCampaign", mock.Anything).Return(nil)

//     req := httptest.NewRequest("POST", "/campaigns", strings.NewReader({"name":"Test"})
// 	)
//     req.Header.Set("Content-Type", "application/json")
//     w := httptest.NewRecorder()

//     h := handler.NewHandler(mockService, nil, nil) // assume handler takes campaign, targeting, delivery
//     h.CreateCampaign(w, req)

//     res := w.Result()
//     assert.Equal(t, http.StatusOK, res.StatusCode)
//     mockService.AssertExpectations(t)
// }

// // test/delivery_test.go
// package test

// import (
//     "net/http"
//     "net/http/httptest"
//     "testing"

//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/mock"
//     "targeting-engine/internal/handler"
// )

// type MockDeliveryService struct {
//     mock.Mock
// }

// func (m *MockDeliveryService) GetDelivery(params map[string]string) (map[string]interface{}, error) {
//     args := m.Called(params)
//     return args.Get(0).(map[string]interface{}), args.Error(1)
// }

// func TestDeliveryAPIWithMock(t *testing.T) {
//     mockService := new(MockDeliveryService)
//     mockService.On("GetDelivery", mock.Anything).Return(map[string]interface{}{"campaign": "test"}, nil)

//     req := httptest.NewRequest("GET", "/v1/delivery?app=news&country=IN&os=android", nil)
//     w := httptest.NewRecorder()

//     h := handler.NewHandler(nil, nil, mockService)
//     h.GetDelivery(w, req)

//     res := w.Result()
//     assert.Equal(t, http.StatusOK, res.StatusCode)
//     mockService.AssertExpectations(t)
// }
