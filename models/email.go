package models

type EmailRequest struct {
	From     string   `json:"from" validate:"required,email"`
	To       []string `json:"to" validate:"required,min=1,dive,email"`
	Subject  string   `json:"subject" validate:"required"`
	HtmlBody string   `json:"htmlBody" validate:"required_without=TextBody"`
	TextBody string   `json:"textBody" validate:"required_without=HtmlBody"`
}

type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
