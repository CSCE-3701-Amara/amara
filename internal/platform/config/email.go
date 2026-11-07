package config

import "fmt"

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func newEmailConfig(lookup func(string) (string, bool)) (EmailConfig, error) {
	value := func(name string) string {
		raw, _ := lookup(name)
		return raw
	}
	cfg := EmailConfig{
		SMTPHost: value("SMTP_HOST"), SMTPPort: value("SMTP_PORT"),
		SMTPUsername: value("SMTP_USERNAME"), SMTPPassword: value("SMTP_PASSWORD"),
		SMTPFrom: value("SMTP_FROM"),
	}
	for name, field := range map[string]string{
		"SMTP_HOST": cfg.SMTPHost, "SMTP_PORT": cfg.SMTPPort, "SMTP_FROM": cfg.SMTPFrom,
	} {
		if field == "" {
			return EmailConfig{}, fmt.Errorf("%s must not be empty", name)
		}
	}
	if (cfg.SMTPUsername == "") != (cfg.SMTPPassword == "") {
		return EmailConfig{}, fmt.Errorf("SMTP_USERNAME and SMTP_PASSWORD must be provided together")
	}
	return cfg, nil
}
