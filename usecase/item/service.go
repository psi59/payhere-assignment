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

func (s *Service) Delete(c context.Context, input *DeleteInput) error {
	// 1. 파라메터 체크
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case valid.IsNil(input):
		return domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return errors.WithStack(err)
	}

	user := input.User

	// 2. 아이템 조회
	item, err := s.itemRepository.Get(c, user.ID, input.ItemID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 3. 아이템 삭제
	if err := s.itemRepository.Delete(c, item.UserID, item.ID); err != nil {
		return errors.WithStack(err)
	}

	// 3. 결과 반환
	return nil
}

func (s *Service) Update(c context.Context, input *UpdateInput) error {
	// 1. 파라메터 체크
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case valid.IsNil(input):
		return domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return errors.WithStack(err)
	}

	user := input.User

	// 2. 아이템 조회
	item, err := s.itemRepository.Get(c, user.ID, input.ItemID)
	if err != nil {
		return errors.WithStack(err)
	}

	// 3. 아이템 수정
	param := &repository.UpdateItemInput{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Cost:        input.Cost,
		Category:    input.Category,
		Barcode:     input.Barcode,
		Size:        input.Size,
		ExpiryAt:    input.ExpiryAt,
	}
	if err := s.itemRepository.Update(c, item.UserID, item.ID, param); err != nil {
		return errors.WithStack(err)
	}

	// 3. 결과 반환
	return nil
}

func (s *Service) Find(c context.Context, input *FindInput) (*FindOutput, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	param := &repository.FindItemInput{
		UserID:      input.User.ID,
		Keyword:     input.Keyword,
		SearchAfter: input.SearchAfter,
	}
	findItemOutput, err := s.itemRepository.Find(c, param)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &FindOutput{
		TotalCount:  findItemOutput.TotalCount,
		Items:       findItemOutput.Items,
		HasNext:     findItemOutput.HasNext,
		SearchAfter: findItemOutput.SearchAfter,
	}, nil
}
