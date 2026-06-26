package main

import (
	"fmt"
	"bankapi/pkg/smtp"
	"bankapi/internal/config"
)

func main() {
	cfg := &config.Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: "587",
		SMTPUser: "test@example.com",
		SMTPPass: "testpass",
	}
	err := smtp.SendEmail("user@example.com", "Test Subject", "<h1>Test Body</h1>", cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("SMTP Test Success (Mocked log should appear)")
	}
}
