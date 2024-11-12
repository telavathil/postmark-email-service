package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/keighl/postmark"
)

func TestMain(m *testing.M) {
	// Set up test environment
	if os.Getenv("POSTMARK_SERVER_TOKEN") == "" {
		os.Setenv("POSTMARK_SERVER_TOKEN", "test-token")
	}

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

type mockPostmarkClient struct{}

func (m *mockPostmarkClient) SendEmail(email postmark.Email) (postmark.EmailResponse, error) {
	return postmark.EmailResponse{}, nil
}

func setupTestApp() *iris.Application {
	app := iris.New()

	// Use mock client instead of real Postmark client
	client := &mockPostmarkClient{}

	// Health check endpoint
	app.Get("/health", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"status": "healthy",
		})
	})

	// Email sending endpoint
	app.Post("/send-email", func(ctx iris.Context) {
		var req EmailRequest
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(EmailResponse{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		// Validate recipients
		if len(req.To) == 0 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(EmailResponse{
				Success: false,
				Message: "At least one recipient is required",
			})
			return
		}

		// Create Postmark email
		email := postmark.Email{
			From:     req.From,
			To:       req.To[0],
			Subject:  req.Subject,
			HtmlBody: req.HtmlBody,
			TextBody: req.TextBody,
		}

		// Send email
		_, err := client.SendEmail(email)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(EmailResponse{
				Success: false,
				Message: "Failed to send email: " + err.Error(),
			})
			return
		}

		ctx.JSON(EmailResponse{
			Success: true,
			Message: "Email sent successfully",
		})
	})

	return app
}

func TestHealthEndpoint(t *testing.T) {
	app := setupTestApp()
	e := httptest.New(t, app)

	// Test health endpoint
	e.GET("/health").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("status").
		HasValue("status", "healthy")
}

func TestSendEmailEndpoint(t *testing.T) {
	app := setupTestApp()
	e := httptest.New(t, app)

	tests := []struct {
		name           string
		request        EmailRequest
		expectedStatus int
		expectedBody   EmailResponse
	}{
		{
			name: "Valid email request",
			request: EmailRequest{
				From:     "test@example.com",
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: http.StatusOK,
			expectedBody: EmailResponse{
				Success: true,
				Message: "Email sent successfully",
			},
		},
		{
			name: "Empty recipient",
			request: EmailRequest{
				From:     "test@example.com",
				To:       []string{},
				Subject:  "Test Subject",
				HtmlBody: "<p>Test body</p>",
				TextBody: "Test body",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: EmailResponse{
				Success: false,
				Message: "At least one recipient is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.request)

			e.POST("/send-email").
				WithBytes(requestBody).
				WithHeader("Content-Type", "application/json").
				Expect().
				Status(tt.expectedStatus).
				JSON().Object().
				ContainsKey("success").
				ContainsKey("message").
				HasValue("success", tt.expectedBody.Success).
				HasValue("message", tt.expectedBody.Message)
		})
	}
}

func TestInvalidJSONRequest(t *testing.T) {
	app := setupTestApp()
	e := httptest.New(t, app)

	// Test invalid JSON request
	invalidJSON := []byte(`{"invalid json"}`)

	e.POST("/send-email").
		WithBytes(invalidJSON).
		WithHeader("Content-Type", "application/json").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("success").
		ContainsKey("message").
		HasValue("success", false).
		HasValue("message", "Invalid request body")
}
