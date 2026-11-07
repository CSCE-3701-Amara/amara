package email

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPSender struct {
	config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) (*SMTPSender, error) {
	for name, value := range map[string]string{
		"host": config.Host, "port": config.Port, "from": config.From,
	} {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("SMTP %s must not be empty", name)
		}
	}
	if strings.ContainsAny(config.From, "\r\n") {
		return nil, errors.New("SMTP from contains invalid characters")
	}
	if address, err := mail.ParseAddress(config.From); err != nil || address.Address != config.From {
		return nil, errors.New("SMTP from must be a valid email address")
	}
	if (config.Username == "") != (config.Password == "") {
		return nil, errors.New("SMTP username and password must be provided together")
	}
	return &SMTPSender{config: config}, nil
}

func (s *SMTPSender) SendVerificationCode(ctx context.Context, recipient, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	recipient = strings.TrimSpace(recipient)
	if recipient == "" {
		return errors.New("verification recipient must not be empty")
	}
	if strings.ContainsAny(recipient, "\r\n") {
		return errors.New("verification recipient contains invalid characters")
	}
	if address, err := mail.ParseAddress(recipient); err != nil || address.Address != recipient {
		return errors.New("verification recipient must be a valid email address")
	}
	if !isVerificationCode(code) {
		return errors.New("verification code must contain only hexadecimal characters")
	}
	auth := smtp.Auth(nil)
	if s.config.Username != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}
	message := buildVerificationMessage(s.config.From, recipient, code)
	return smtp.SendMail(
		fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		auth,
		s.config.From,
		[]string{recipient},
		message,
	)
}

func isVerificationCode(code string) bool {
	if code == "" {
		return false
	}
	for _, character := range code {
		if !((character >= '0' && character <= '9') ||
			(character >= 'a' && character <= 'f') ||
			(character >= 'A' && character <= 'F')) {
			return false
		}
	}
	return true
}

func buildVerificationMessage(from, recipient, code string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: AMARA email verification\r\n\r\nYour AMARA verification code is: %s\r\n",
		from, recipient, code,
	))
}
