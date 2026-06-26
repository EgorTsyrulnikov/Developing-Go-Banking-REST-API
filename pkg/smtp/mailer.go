package smtp

import (
	"crypto/tls"
	"log"

	"github.com/go-mail/mail/v2"
	"bankapi/internal/config"
	"strconv"
)

func SendEmail(to string, subject string, body string, cfg *config.Config) error {
	port, _ := strconv.Atoi(cfg.SMTPPort)

	m := mail.NewMessage()
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPass)
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.SMTPHost,
		InsecureSkipVerify: true, // For mock testing
	}

	// Just log it since it's mock
	log.Printf("MOCK SMTP: Sending email to %s, Subject: %s", to, subject)
	
	// If you want to actually connect, uncomment this:
	// if err := d.DialAndSend(m); err != nil {
	// 	return fmt.Errorf("email sending failed: %v", err)
	// }
	
	return nil
}
