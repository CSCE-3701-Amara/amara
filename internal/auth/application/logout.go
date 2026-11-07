package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type LogoutService struct {
	sessions ports.SessionTerminator
}

func NewLogoutService(sessions ports.SessionTerminator) (*LogoutService, error) {
	if sessions == nil {
		return nil, errors.New("session terminator must not be nil")
	}
	return &LogoutService{sessions: sessions}, nil
}

func (s *LogoutService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrInvalidSession
	}
	if err := s.sessions.TerminateSession(ctx, token); err != nil {
		return fmt.Errorf("terminate session: %w", err)
	}
	return nil
}
