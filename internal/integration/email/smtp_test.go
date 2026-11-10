package email

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestBuildSMTPHeadersKeepsRecipientOutOfMessage(t *testing.T) {
	date := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	headers := buildSMTPHeaders(&mail.Address{Address: "no-reply@example.com"}, date)
	want := "From: <no-reply@example.com>\r\nTo: undisclosed-recipients:;\r\nDate: Tue, 10 Nov 2026 12:00:00 +0000\r\n"
	if headers != want {
		t.Fatalf("headers = %q, want %q", headers, want)
	}
	message, err := mail.ReadMessage(strings.NewReader(headers + string(buildVerificationMessage("deadbeef"))))
	if err != nil {
		t.Fatal(err)
	}
	addresses, err := message.Header.AddressList("To")
	if err != nil || len(addresses) != 0 {
		t.Fatalf("To must be a valid empty group; got %v, error = %v", addresses, err)
	}
	if len(message.Header) != 4 || message.Header.Get("Bcc") != "" {
		t.Fatalf("unexpected message headers: %v", message.Header)
	}
}

func TestSMTPSenderRejectsUntrustedRecipientBeforeConnecting(t *testing.T) {
	sender, err := NewSMTPSender(SMTPConfig{Host: "127.0.0.1", Port: "0", From: "no-reply@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	for _, recipient := range []string{
		"", "   ", "not-an-email",
		"user@example.com\r\nBcc: victim@example.com",
		"user@example.com\nBcc: victim@example.com",
		"user@example.com\rBcc: victim@example.com",
		"user@example.com\r\n\r\nInjected message body",
		"user@example.com, victim@example.com",
		"Display Name <user@example.com>",
	} {
		t.Run(recipient, func(t *testing.T) {
			for name, send := range map[string]func() error{
				"verification": func() error { return sender.SendVerificationCode(context.Background(), recipient, "deadbeef") },
				"MFA":          func() error { return sender.SendLoginCode(context.Background(), recipient, "123456") },
			} {
				t.Run(name, func(t *testing.T) {
					if err := send(); err == nil || !strings.Contains(err.Error(), "recipient") {
						t.Fatalf("invalid recipient must fail before connecting to SMTP; got %v", err)
					}
				})
			}
		})
	}
}

func TestNewSMTPSenderRejectsHeaderInjectionInFrom(t *testing.T) {
	_, err := NewSMTPSender(SMTPConfig{
		Host: "smtp.example.com",
		Port: "587",
		From: "no-reply@example.com\r\nBcc: victim@example.com",
	})
	if err == nil {
		t.Fatal("NewSMTPSender() accepted a From value containing a newline")
	}
}

func TestNewSMTPSenderRejectsInvalidFrom(t *testing.T) {
	_, err := NewSMTPSender(SMTPConfig{
		Host: "smtp.example.com",
		Port: "587",
		From: "not-an-email",
	})
	if err == nil {
		t.Fatal("NewSMTPSender() accepted an invalid From value")
	}
}

func TestIsVerificationCode(t *testing.T) {
	for _, test := range []struct {
		name string
		code string
		want bool
	}{
		{name: "hexadecimal", code: "a1B20f", want: true},
		{name: "empty", code: "", want: false},
		{name: "header injection", code: "abc\r\nBcc: victim@example.com", want: false},
		{name: "html", code: "<script>", want: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := isVerificationCode(test.code); got != test.want {
				t.Fatalf("isVerificationCode(%q) = %t, want %t", test.code, got, test.want)
			}
		})
	}
}

func TestBuildVerificationMessage(t *testing.T) {
	message := string(buildVerificationMessage("deadbeef"))
	want := "Subject: AMARA email verification\r\n\r\n" +
		"Your AMARA verification code is: deadbeef\r\n"

	if message != want {
		t.Fatalf("buildVerificationMessage() = %q, want %q", message, want)
	}
}

func TestSendVerificationCodeRequiresSTARTTLS(t *testing.T) {
	listener := listenSMTPTestRelay(t, func(conn net.Conn) {
		defer conn.Close()
		_, _ = fmt.Fprint(conn, "220 test relay\r\n")
		reader := bufio.NewReader(conn)
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if !strings.HasPrefix(line, "EHLO ") {
			return
		}
		_, _ = fmt.Fprint(conn, "250-test relay\r\n250 AUTH PLAIN\r\n")
		remaining, _ := reader.ReadString('\n')
		if strings.Contains(remaining, "deadbeef") {
			t.Errorf("relay received verification code without STARTTLS")
		}
	})
	sender := testSMTPSender(t, listener, 2*time.Second)

	err := sender.SendVerificationCode(context.Background(), "user@example.com", "deadbeef")
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Fatalf("SendVerificationCode() error = %v, want STARTTLS error", err)
	}
}

func TestSendVerificationCodeStopsWhenContextIsCancelled(t *testing.T) {
	listener := listenSMTPTestRelay(t, func(conn net.Conn) {
		defer conn.Close()
		time.Sleep(time.Second)
	})
	sender := testSMTPSender(t, listener, 10*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := sender.SendVerificationCode(ctx, "user@example.com", "deadbeef")
	if err == nil {
		t.Fatal("SendVerificationCode() succeeded against a stalled relay")
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatalf("SendVerificationCode() took too long after cancellation: %s", time.Since(start))
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("context error = %v, want deadline exceeded", ctx.Err())
	}
}

func listenSMTPTestRelay(t *testing.T, handler func(net.Conn)) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			handler(conn)
		}
	}()
	return listener
}

func testSMTPSender(t *testing.T, listener net.Listener, timeout time.Duration) *SMTPSender {
	t.Helper()
	sender, err := NewSMTPSender(SMTPConfig{
		Host: "127.0.0.1", Port: strings.TrimPrefix(listener.Addr().String(), "127.0.0.1:"),
		From: "no-reply@example.com", Timeout: timeout,
	})
	if err != nil {
		t.Fatal(err)
	}
	return sender
}
