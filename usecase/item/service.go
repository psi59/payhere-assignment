package item

import (
	"context"

	"github.com/pkg/errors"

	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/psi59/payhere-assignment/repository"
)

type Service struct {
	itemRepository repository.ItemRepository
}

func NewService(itemRepository repository.ItemRepository) (*Service, error) {
	if valid.IsNil(itemRepository) {
		return nil, repository.ErrNilItemRepository
	}

	return &Service{itemRepository: itemRepository}, nil
}

func (s *Service) Create(c context.Context, input *CreateInput) (*CreateOutput, error) {
	// 1. 파라메터 체크
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	// 2. 도메인 객체 생성
	user := input.User
	item, err := domain.NewItem(
		user.ID,
		input.Name,
		input.Description,
		input.Price,
		input.Cost,
		input.Category,
		input.Barcode,
		input.ExpiryAt,
		input.Size,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 3. 아이템 생성
	if err := s.itemRepository.Create(c, item); err != nil {
		return nil, errors.WithStack(err)
	}

	// 4. 결과 반환
	return &CreateOutput{
		Item: item,
	}, nil
}

func (s *Service) Get(c context.Context, input *GetInput) (*GetOutput, error) {
	// 1. 파라메터 체크
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	user := input.User

	// 2. 아이템 조회
	item, err := s.itemRepository.Get(c, user.ID, input.ItemID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 3. 결과 반환
	return &GetOutput{Item: item}, nil
}
