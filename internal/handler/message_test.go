package handler

import (
	"multi-tenant-messaging-app/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFetchMessages(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMessageServiceInterface(ctrl)
	handler := NewMessageHandler(mockService)

	// Test cases
	testCases := []struct {
		name           string
		tenantID       string
		cursor         string
		limit          int
		mockMessages   []map[string]interface{}
		mockNextCursor string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success",
			tenantID: "test-tenant",
			cursor:   "550e8400-e29b-41d4-a716-446655440000",
			limit:    10,
			mockMessages: []map[string]interface{}{
				{"id": "550e8400-e29b-41d4-a716-446655440001", "payload": map[string]interface{}{"key": "value"}, "created_at": "2023-06-21T15:05:00.000Z"},
				{"id": "550e8400-e29b-41d4-a716-446655440002", "payload": map[string]interface{}{"key": "another value"}, "created_at": "2023-06-21T15:06:00.000Z"},
			},
			mockNextCursor: "550e8400-e29b-41d4-a716-446655440002",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":[{"id":"550e8400-e29b-41d4-a716-446655440001","payload":{"key":"value"},"created_at":"2023-06-21T15:05:00.000Z"},{"id":"550e8400-e29b-41d4-a716-446655440002","payload":{"key":"another value"},"created_at":"2023-06-21T15:06:00.000Z"}],"cursor":"550e8400-e29b-41d4-a716-446655440002"}`,
		},
		// {
		// 	name:           "Internal Server Error",
		// 	tenantID:       "test-tenant",
		// 	cursor:         "550e8400-e29b-41d4-a716-446655440000",
		// 	limit:          10,
		// 	mockMessages:   nil,
		// 	mockNextCursor: "",
		// 	mockError:      errors.New("internal error"),
		// 	expectedStatus: http.StatusInternalServerError,
		// 	expectedBody:   `{"error":"internal error"}`,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock service behavior
			mockService.EXPECT().
				FetchMessages("", tc.cursor, tc.limit).
				Return(tc.mockMessages, tc.mockNextCursor, tc.mockError)

			// Create request and response recorder
			req := httptest.NewRequest(http.MethodGet, "/api/messages?tenant_id="+tc.tenantID+"&cursor="+tc.cursor, nil)
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.FetchMessages(c)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
