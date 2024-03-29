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

func (s *Service) Create(c context.Context, input *CreateInput) (*CreateOutput, error) {
	if valid.IsNil(c) {
		return nil, domain.ErrNilContext
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := domain.NewUser(input.PhoneNumber, input.Password, time.Now())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.userRepository.Create(c, user); err != nil {
		return nil, errors.WithStack(err)
	}

	return &CreateOutput{User: user}, nil
}

func (s *Service) Get(c context.Context, input *GetInput) (*GetOutput, error) {
	if valid.IsNil(c) {
		return nil, domain.ErrNilContext
	}
	if valid.IsNil(input) {
		return nil, domain.ErrNilInput
	}
	if err := valid.ValidateStruct(input); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := s.userRepository.Get(c, input.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &GetOutput{User: user}, nil
}

func (s *Service) GetByPhoneNumber(c context.Context, input *GetByPhoneNumberInput) (*GetOutput, error) {
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

	return &GetOutput{User: user}, nil
}
