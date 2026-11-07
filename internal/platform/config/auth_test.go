package config

import (
	"strings"
	"testing"
)

func TestLoadPreservesMFAEncryptionKeyBytes(t *testing.T) {
	for _, key := range []string{
		"0123456789abcdef0123456789abcdef",
		" " + strings.Repeat("k", 30) + " ",
		strings.Repeat("é", 16),
	} {
		cfg, err := fromValues(map[string]string{"AUTH_MFA_ENCRYPTION_KEY": key})
		if err != nil {
			t.Fatal(err)
		}
		if string(cfg.Auth.MFAEncryptionKey[:]) != key {
			t.Fatal("configuration must preserve the exact MFA key bytes")
		}
	}
}

func TestLoadRejectsInvalidMFAEncryptionKey(t *testing.T) {
	for _, key := range []string{"", "sensitive-short-key", strings.Repeat("k", 31), strings.Repeat("k", 33)} {
		_, err := fromValues(map[string]string{"AUTH_MFA_ENCRYPTION_KEY": key})
		if err == nil || err.Error() != "AUTH_MFA_ENCRYPTION_KEY must be exactly 32 bytes" {
			t.Fatalf("expected key validation without echoing the secret, got %v", err)
		}
	}
}

func TestLoadRequiresMFAEncryptionKey(t *testing.T) {
	_, err := load(func(name string) (string, bool) {
		if name == "DATABASE_URL" {
			return "postgres://user:pass@localhost:5432/db", true
		}
		return "", false
	})
	if err == nil || !strings.Contains(err.Error(), "AUTH_MFA_ENCRYPTION_KEY") {
		t.Fatalf("expected missing MFA key error, got %v", err)
	}
}
