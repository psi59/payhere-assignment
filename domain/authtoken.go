package domain

import (
	"fmt"
	"time"
)

type AuthToken struct {
	Token     string
	ExpiresAt time.Time
}

const ErrNilAuthToken ConstantError = "nil AuthToken"

func (t *AuthToken) Validate() error {
	if len(t.Token) == 0 {
		return fmt.Errorf("empty token")
	}
	if t.ExpiresAt.IsZero() {
		return fmt.Errorf("zero expiresAt")
	}

	return nil
}
