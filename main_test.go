package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"postmark-email-service/config"
	"postmark-email-service/models"
	"postmark-email-service/server"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/keighl/postmark"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up test environment
	if os.Getenv("POSTMARK_TOKEN") == "" {
		os.Setenv("POSTMARK_TOKEN", "test-token")
	}

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

type mockPostmarkClient struct{}

func (m *mockPostmarkClient) SendEmail(email postmark.Email) (postmark.EmailResponse, error) {
	return postmark.EmailResponse{
		ErrorCode: 0,
		Message:   "OK",
		To:        email.To,
	}, nil
}

func setupTestServer() *server.Server {
	cfg := &config.Config{
		PostmarkToken: "test-token",
		Port:         "8080",
	}
	srv := server.NewServer(cfg.PostmarkToken)
	srv.Client = &mockPostmarkClient{} // Use mock client
	return srv
}

func TestHTTPEndpoint(t *testing.T) {
	srv := setupTestServer()

	tests := []struct {
		name           string
		method         string
		request        models.EmailRequest
		expectedStatus int
		expectedBody   models.EmailResponse
	}{
		{
			name:   "Valid email request",
			method: http.MethodPost,
			request: models.EmailRequest{
				From:     "test@example.com",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: http.StatusOK,
			expectedBody: models.EmailResponse{
				Success: true,
				Message: "Email sent successfully",
			},
		},
		{
			name:   "Empty recipient",
			method: http.MethodPost,
			request: models.EmailRequest{
				From:     "test@example.com",
				To:       []string{},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.EmailResponse{
				Success: false,
				Message: "Key: 'EmailRequest.To' Error:Field validation for 'To' failed on the 'min' tag",
			},
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			request:        models.EmailRequest{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody: models.EmailResponse{
				Success: false,
				Message: "Method not allowed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.method == http.MethodPost {
				body, err = json.Marshal(tt.request)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(tt.method, "/send-email", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			srv.HandleHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.EmailResponse
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody.Success, response.Success)
			assert.Equal(t, tt.expectedBody.Message, response.Message)
		})
	}
}

func TestLambdaHandler(t *testing.T) {
	srv := setupTestServer()

	tests := []struct {
		name           string
		request        models.EmailRequest
		expectedStatus int
		expectedBody   models.EmailResponse
	}{
		{
			name: "Valid email request",
			request: models.EmailRequest{
				From:     "test@example.com",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: 200,
			expectedBody: models.EmailResponse{
				Success: true,
				Message: "Email sent successfully",
			},
		},
		{
			name: "Empty recipient",
			request: models.EmailRequest{
				From:     "test@example.com",
				To:       []string{},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: 400,
			expectedBody: models.EmailResponse{
				Success: false,
				Message: "Key: 'EmailRequest.To' Error:Field validation for 'To' failed on the 'min' tag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			event := events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:      string(body),
			}

			response, err := srv.HandleRequest(event)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			var responseBody models.EmailResponse
			err = json.Unmarshal([]byte(response.Body), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody.Success, responseBody.Success)
			assert.Equal(t, tt.expectedBody.Message, responseBody.Message)
		})
	}
}
