package domain

import (
	"errors"
	"fmt"
)

type AccountAction string

const (
	AccountActionDeactivate AccountAction = "deactivate"
	AccountActionDelete     AccountAction = "delete"
)

var ErrInvalidAccountAction = errors.New("invalid account action")

func (a AccountAction) Validate() error {
	if a != AccountActionDeactivate && a != AccountActionDelete {
		return fmt.Errorf("%w: action is invalid", ErrInvalidAccountAction)
	}
	return nil
}
