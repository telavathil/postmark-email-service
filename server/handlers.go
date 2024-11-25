package server

import (
	"encoding/json"
	"net/http"
	"postmark-email-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/keighl/postmark"
)

// HandleHTTP handles HTTP requests
func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.EmailResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req models.EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.EmailResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := s.Validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.EmailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	email := &postmark.Email{
		From:     req.From,
		To:       req.To[0],
		Subject:  req.Subject,
		HtmlBody: req.HtmlBody,
		TextBody: req.TextBody,
	}

	_, err := s.Client.SendEmail(*email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.EmailResponse{
			Success: false,
			Message: "Failed to send email: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.EmailResponse{
		Success: true,
		Message: "Email sent successfully",
	})
}

// HandleRequest handles Lambda requests
func (s *Server) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.EmailRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{\"success\": false, \"message\": \"Invalid request body\"}",
		}, nil
	}

	if err := s.Validate.Struct(req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{\"success\": false, \"message\": \"" + err.Error() + "\"}",
		}, nil
	}

	email := &postmark.Email{
		From:     req.From,
		To:       req.To[0],
		Subject:  req.Subject,
		HtmlBody: req.HtmlBody,
		TextBody: req.TextBody,
	}

	_, err := s.Client.SendEmail(*email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{\"success\": false, \"message\": \"Failed to send email: " + err.Error() + "\"}",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "{\"success\": true, \"message\": \"Email sent successfully\"}",
	}, nil
}
