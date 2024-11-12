package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/keighl/postmark"
)

type Server struct {
	Validate *validator.Validate
	Client   *postmark.Client
}

func NewServer(postmarkToken string) *Server {
	return &Server{
		Validate: validator.New(),
		Client:   postmark.NewClient(postmarkToken, ""),
	}
}
