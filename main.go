package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"postmark-email-service/config"
	"postmark-email-service/models"
	"postmark-email-service/server"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keighl/postmark"
)

func main() {
	cfg := config.LoadConfig()
	server := server.NewServer(cfg.PostmarkToken)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		lambda.Start(server.HandleRequest)
	} else {
		// Running locally
		http.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				server.ErrorResponse(http.StatusMethodNotAllowed, "Method not allowed")
				return
			}

			var req models.EmailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				server.ErrorResponse(http.StatusBadRequest, "Invalid request body")
				return
			}

			if err := server.Validate.Struct(req); err != nil {
				server.ErrorResponse(http.StatusBadRequest, err.Error())
				return
			}

			email := &postmark.Email{
				From:     req.From,
				To:       req.To[0],
				Subject:  req.Subject,
				HtmlBody: req.HtmlBody,
				TextBody: req.TextBody,
			}

			_, err := server.Client.SendEmail(*email)
			if err != nil {
				server.ErrorResponse(http.StatusInternalServerError, "Failed to send email: "+err.Error())
				return
			}

			server.SuccessResponse()
		})

		log.Printf("Server starting on :%s", cfg.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
	}
}
