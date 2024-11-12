package server

import (
	"encoding/json"
	"net/http"
	"postmark-service/models"

	"github.com/aws/aws-lambda-go/events"
)

func (s *Server) ErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body: string(mustJSON(models.EmailResponse{
			Success: false,
			Message: message,
		})),
	}, nil
}

func (s *Server) SuccessResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(mustJSON(models.EmailResponse{
			Success: true,
			Message: "Email sent successfully",
		})),
	}, nil
}

func mustJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
