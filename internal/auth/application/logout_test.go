package application

import (
	"context"
	"errors"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type logoutSessions struct {
	token string
	err   error
}

func (s *logoutSessions) TerminateSession(_ context.Context, token string) error {
	s.token = token
	return s.err
}

func TestLogoutTerminatesSession(t *testing.T) {
	sessions := &logoutSessions{}
	service, err := NewLogoutService(sessions)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Logout(context.Background(), "session-token"); err != nil {
		t.Fatal(err)
	}
	if sessions.token != "session-token" {
		t.Fatalf("terminated token = %q, want session-token", sessions.token)
	}
}

func TestLogoutRejectsEmptyTokenAndMapsPersistenceError(t *testing.T) {
	sessions := &logoutSessions{err: errors.New("storage down")}
	service, err := NewLogoutService(sessions)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Logout(context.Background(), ""); !errors.Is(err, domain.ErrInvalidSession) {
		t.Fatalf("empty token error = %v, want invalid session", err)
	}
	if err := service.Logout(context.Background(), "session-token"); err == nil {
		t.Fatal("expected persistence error")
	}
}
