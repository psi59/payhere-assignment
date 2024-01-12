package user

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/psi59/payhere-assignment/repository"
)

type Service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) (*Service, error) {
	if valid.IsNil(userRepository) {
		return nil, repository.ErrNilUserRepository
	}

	return &Service{
		userRepository: userRepository,
	}, nil
}

func (s *Service) Create(c context.Context, input *CreateInput) (*domain.User, error) {
	if valid.IsNil(c) {
		return nil, domain.ErrNilContext
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := domain.NewUser(input.Name, input.PhoneNumber, input.Password, time.Now())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.userRepository.Create(c, user); err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func (s *Service) GetByPhoneNumber(c context.Context, input *GetByPhoneNumberInput) (*domain.User, error) {
	if valid.IsNil(c) {
		return nil, domain.ErrNilContext
	}
	if valid.IsNil(input) {
		return nil, domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := s.userRepository.GetByPhoneNumber(c, input.PhoneNumber)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}
