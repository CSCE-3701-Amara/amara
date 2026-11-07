package config

import "fmt"

// AuthConfig holds authentication settings supplied by the environment.
type AuthConfig struct {
	// MFAEncryptionKey encrypts MFA secrets at rest and must never be logged.
	MFAEncryptionKey [32]byte
}

func newAuthConfig(lookup func(string) (string, bool)) (AuthConfig, error) {
	var cfg AuthConfig
	// Read the raw bytes: trimming or decoding would change an existing key.
	// The deployment supplies the key; there is no default.
	raw, _ := lookup("AUTH_MFA_ENCRYPTION_KEY")
	if len(raw) != len(cfg.MFAEncryptionKey) {
		return AuthConfig{}, fmt.Errorf("AUTH_MFA_ENCRYPTION_KEY must be exactly 32 bytes")
	}
	copy(cfg.MFAEncryptionKey[:], raw)
	return cfg, nil
}
