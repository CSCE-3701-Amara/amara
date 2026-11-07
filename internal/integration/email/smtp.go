package email

import (
	"context"
	"errors"
	"fmt"
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
	auth := smtp.Auth(nil)
	if s.config.Username != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}
	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: AMARA email verification\r\n\r\nYour AMARA verification code is: %s\r\n",
		s.config.From, recipient, code,
	))
	return smtp.SendMail(
		fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		auth,
		s.config.From,
		[]string{recipient},
		message,
	)
}
