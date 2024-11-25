package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/keighl/postmark"
)

type PostmarkClient interface {
	SendEmail(email postmark.Email) (postmark.EmailResponse, error)
}

type Server struct {
	Client    PostmarkClient
	Validate  *validator.Validate
}

func NewServer(postmarkToken string) *Server {
	return &Server{
		Client:    postmark.NewClient(postmarkToken, ""),
		Validate:  validator.New(),
	}
}
