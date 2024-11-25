package main

import (
	"log"
	"net/http"
	"os"
	"postmark-email-service/config"
	"postmark-email-service/server"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	cfg := config.LoadConfig()
	server := server.NewServer(cfg.PostmarkToken)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		lambda.Start(server.HandleRequest)
	} else {
		// Running locally
		http.HandleFunc("/send-email", server.HandleHTTP)
		log.Printf("Server starting on :%s", cfg.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
	}
}
