package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	Timeout  time.Duration
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
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
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
	message := buildVerificationMessage(code)
	address := net.JoinHostPort(s.config.Host, s.config.Port)
	dialer := net.Dialer{Timeout: s.config.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to SMTP relay: %w", err)
	}
	done := make(chan struct{})
	defer close(done)
	go closeOnContextCancellation(ctx, conn, done)
	if err := conn.SetDeadline(time.Now().Add(s.config.Timeout)); err != nil {
		_ = conn.Close()
		return fmt.Errorf("set SMTP deadline: %w", err)
	}
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()
	if err := client.Hello(s.config.Host); err != nil {
		return fmt.Errorf("SMTP hello: %w", err)
	}
	if ok, _ := client.Extension("STARTTLS"); !ok {
		return errors.New("SMTP relay does not advertise STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{
		ServerName: s.config.Host,
		MinVersion: tls.VersionTLS12,
	}); err != nil {
		return fmt.Errorf("start SMTP TLS: %w", err)
	}
	if state, ok := client.TLSConnectionState(); !ok || !state.HandshakeComplete {
		return errors.New("SMTP TLS handshake was not completed")
	}
	if s.config.Username != "" {
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authenticate with SMTP relay: %w", err)
		}
	}
	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(recipient); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start SMTP message: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("close SMTP session: %w", err)
	}
	return nil
}

func closeOnContextCancellation(ctx context.Context, conn net.Conn, done <-chan struct{}) {
	select {
	case <-ctx.Done():
		_ = conn.Close()
	case <-done:
	}
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

func buildVerificationMessage(code string) []byte {
	return []byte(fmt.Sprintf(
		"Subject: AMARA email verification\r\n\r\nYour AMARA verification code is: %s\r\n",
		html.EscapeString(code),
	))
}
