package domain

import (
	"fmt"
	"time"

	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
}

const (
	ErrNilUser     ConstantError = "nil User"
	ErrInvalidUser ConstantError = "InvalidUser"
)

func NewUser(phoneNumber, password string, createdAt time.Time) (*User, error) {
	if err := valid.ValidatePassword(password); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := valid.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, errors.WithStack(err)
	}
	if createdAt.IsZero() {
		return nil, fmt.Errorf("zero createdAt")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate hashed password")
	}
	u := &User{
		ID:          0,
		PhoneNumber: phoneNumber,
		Password:    string(hashed),
		CreatedAt:   createdAt,
	}

	return u, nil
}

func (u *User) Validate() error {
	if u.ID == 0 {
		return fmt.Errorf("zero UserID")
	}
	if _, err := bcrypt.Cost([]byte(u.Password)); err != nil {
		return errors.Wrap(ErrInvalidUser, err.Error())
	}
	if err := valid.ValidatePhoneNumber(u.PhoneNumber); err != nil {
		return errors.Wrap(ErrInvalidUser, err.Error())
	}
	if u.CreatedAt.IsZero() {
		return fmt.Errorf("%w: zero time", ErrInvalidUser)
	}

	return nil
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.Wrap(ErrPasswordMismatch, err.Error())
	}

	return nil
}
