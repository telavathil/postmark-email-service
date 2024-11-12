package server

import (
	"encoding/json"
	"net/http"
	"postmark-email-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/keighl/postmark"
)

func (s *Server) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.EmailRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return s.ErrorResponse(http.StatusBadRequest, "Invalid request body")
	}

	if err := s.Validate.Struct(req); err != nil {
		return s.ErrorResponse(http.StatusBadRequest, err.Error())
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
		return s.ErrorResponse(http.StatusInternalServerError, "Failed to send email: "+err.Error())
	}

	return s.SuccessResponse()
}
