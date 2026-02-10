package proxy

import (
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"gopkg.in/mail.v2"
)

type EmailService interface {
	Send(data, emailTo, subject string) error
}

type EmailSender struct {
	SmtpHost      string `json:"smtp_host"`
	SmtpEmailFrom string `json:"smtp_email_from"`
	SmtpPass      string `json:"smtp_pass"`
}

var (
	instance *EmailSender
	once     sync.Once
)

// GetEmailInstance returns EmailSender singleton instance
func GetEmailInstance() *EmailSender {
	once.Do(func() {
		eConfig := config.Config.EmailConfig
		instance = &EmailSender{
			SmtpHost:      eConfig.SmtpHost,
			SmtpEmailFrom: eConfig.SmtpEmailFrom,
			SmtpPass:      eConfig.SmtpPass,
		}
	})
	return instance
}

// Send Email
func (s *EmailSender) Send(data, emailTo, subject string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.SmtpEmailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)
	d := mail.NewDialer(s.SmtpHost, 465, s.SmtpEmailFrom, s.SmtpPass)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
