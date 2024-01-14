package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/psi59/payhere-assignment/internal/valid"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(c context.Context, user *domain.User) error {
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case valid.IsNil(user):
		return domain.ErrNilUser
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	userModel := &User{
		UserID:      user.ID,
		UserName:    user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
		CreatedAt:   user.CreatedAt,
	}
	if err := conn.Create(userModel).Error; err != nil {
		if IsDuplicateEntry(err) {
			return errors.Wrap(domain.ErrDuplicatedUser, err.Error())
		}

		return errors.WithStack(err)
	}
	user.ID = userModel.UserID

	return nil
}

func (r *UserRepository) Get(c context.Context, userID int) (*domain.User, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilUser
	case userID < 1:
		return nil, fmt.Errorf("invalid userID")
	}

	conn, err := db.ConnFromContext(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userModel := User{UserID: userID}
	if err := conn.Take(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &domain.User{
		ID:          userModel.UserID,
		Name:        userModel.UserName,
		PhoneNumber: userModel.PhoneNumber,
		Password:    userModel.Password,
		CreatedAt:   userModel.CreatedAt,
	}, nil
}

func (r *UserRepository) GetByPhoneNumber(c context.Context, phoneNumber string) (*domain.User, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilUser
	case len(phoneNumber) == 0:
		return nil, fmt.Errorf("empty phoneNumber")
	}

	conn, err := db.ConnFromContext(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var userModel User
	if err := conn.Where("phone_number = ?", phoneNumber).Take(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &domain.User{
		ID:          userModel.UserID,
		Name:        userModel.UserName,
		PhoneNumber: userModel.PhoneNumber,
		Password:    userModel.Password,
		CreatedAt:   userModel.CreatedAt,
	}, nil
}

type User struct {
	UserID      int       `gorm:"user_id;primaryKey"`
	UserName    string    `gorm:"user_name"`
	PhoneNumber string    `gorm:"phone_number"`
	Password    string    `gorm:"password"`
	CreatedAt   time.Time `gorm:"created_at"`
}

func (u *User) TableName() string {
	return "users"
}
