package email

import "testing"

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
	message := string(buildVerificationMessage(
		"no-reply@example.com",
		"user@example.com",
		"deadbeef",
	))
	want := "From: no-reply@example.com\r\n" +
		"To: user@example.com\r\n" +
		"Subject: AMARA email verification\r\n\r\n" +
		"Your AMARA verification code is: deadbeef\r\n"

	if message != want {
		t.Fatalf("buildVerificationMessage() = %q, want %q", message, want)
	}
}
