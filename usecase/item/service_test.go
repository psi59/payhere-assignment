package item

import (
	"context"
	"testing"
	"time"

	"github.com/psi59/payhere-assignment/repository"

	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/mocks/repomocks"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var userDomain *domain.User

func init() {
	u, err := domain.NewUser(
		gofakeit.Regex(`^01\d{8,9}$`),
		gofakeit.Password(true, true, true, true, true, 10),
		gofakeit.Date(),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	u.ID = gofakeit.Number(1, 10)

	userDomain = u
}

func TestService_Create(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepository := repomocks.NewMockItemRepository(ctrl)
	srv, err := NewService(itemRepository)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		input := &CreateInput{
			User:        userDomain,
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    gofakeit.FutureDate(),
		}
		itemRepository.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, item *domain.Item) error {
			item.ID = gofakeit.Number(1, 10)
			return nil
		})

		got, err := srv.Create(ctx, input)
		assert.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, got.Item.ID > 0)
	})

	t.Run("nil context", func(t *testing.T) {
		input := &CreateInput{
			User:        userDomain,
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    gofakeit.FutureDate(),
		}

		got, err := srv.Create(nil, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		got, err := srv.Create(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		input := &CreateInput{
			User:        nil,
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    gofakeit.FutureDate(),
		}

		got, err := srv.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("중복 아이템 생성 실패", func(t *testing.T) {
		input := &CreateInput{
			User:        userDomain,
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    gofakeit.FutureDate(),
		}
		itemRepository.EXPECT().Create(ctx, gomock.Any()).Return(domain.ErrItemAlreadyExists)

		got, err := srv.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("아이템 생성 실패", func(t *testing.T) {
		input := &CreateInput{
			User:        userDomain,
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    gofakeit.FutureDate(),
		}
		itemRepository.EXPECT().Create(ctx, gomock.Any()).Return(gofakeit.Error())

		got, err := srv.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestService_Get(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepository := repomocks.NewMockItemRepository(ctrl)
	srv, err := NewService(itemRepository)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(item, nil)
		input := &GetInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		got, err := srv.Get(ctx, input)
		assert.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, item, got.Item)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		input := &GetInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		got, err := srv.Get(nil, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		got, err := srv.Get(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		input := &GetInput{
			User:   userDomain,
			ItemID: 0,
		}
		got, err := srv.Get(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("item not found", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(nil, domain.ErrItemNotFound)
		input := &GetInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		got, err := srv.Get(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("아이템 조회 에러", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(nil, gofakeit.Error())
		input := &GetInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		got, err := srv.Get(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepository := repomocks.NewMockItemRepository(ctrl)
	srv, err := NewService(itemRepository)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(item, nil)
		itemRepository.EXPECT().Delete(ctx, userDomain.ID, item.ID).Return(nil)
		input := &DeleteInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		err := srv.Delete(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		input := &DeleteInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		err := srv.Delete(nil, input)
		assert.Error(t, err)
	})

	t.Run("nil input", func(t *testing.T) {
		err := srv.Delete(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		input := &DeleteInput{
			User:   userDomain,
			ItemID: 0,
		}
		err := srv.Delete(ctx, input)
		assert.Error(t, err)
	})

	t.Run("item not found", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(nil, domain.ErrItemNotFound)
		input := &DeleteInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		err := srv.Delete(ctx, input)
		assert.Error(t, err)
	})

	t.Run("아이템 삭제 에러", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(item, nil)
		itemRepository.EXPECT().Delete(ctx, userDomain.ID, item.ID).Return(gofakeit.Error())
		input := &DeleteInput{
			User:   userDomain,
			ItemID: item.ID,
		}
		err := srv.Delete(ctx, input)
		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepository := repomocks.NewMockItemRepository(ctrl)
	srv, err := NewService(itemRepository)
	assert.NoError(t, err)
	item := newTestItem(t, userDomain.ID)

	name := gofakeit.Drink()
	updateInput := &repository.UpdateItemInput{
		Name: &name,
	}

	t.Run("OK", func(t *testing.T) {
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(item, nil)
		itemRepository.EXPECT().Update(ctx, userDomain.ID, item.ID, updateInput).Return(nil)
		input := &UpdateInput{
			User:   userDomain,
			ItemID: item.ID,
			Name:   &name,
		}
		err := srv.Update(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		input := &UpdateInput{
			User:   userDomain,
			ItemID: item.ID,
			Name:   &name,
		}
		err := srv.Update(nil, input)
		assert.Error(t, err)
	})

	t.Run("nil input", func(t *testing.T) {
		err := srv.Update(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		input := &UpdateInput{
			User:   userDomain,
			ItemID: 0,
			Name:   &name,
		}
		err := srv.Update(ctx, input)
		assert.Error(t, err)
	})

	t.Run("item not found", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(nil, domain.ErrItemNotFound)
		input := &UpdateInput{
			User:   userDomain,
			ItemID: item.ID,
			Name:   &name,
		}
		err := srv.Update(ctx, input)
		assert.Error(t, err)
	})

	t.Run("아이템 수정 에러", func(t *testing.T) {
		item := newTestItem(t, userDomain.ID)
		itemRepository.EXPECT().Get(ctx, userDomain.ID, item.ID).Return(item, nil)
		itemRepository.EXPECT().Update(ctx, userDomain.ID, item.ID, updateInput).Return(gofakeit.Error())
		input := &UpdateInput{
			User:   userDomain,
			ItemID: item.ID,
			Name:   &name,
		}
		err := srv.Update(ctx, input)
		assert.Error(t, err)
	})
}

func newTestItem(t *testing.T, userID int) *domain.Item {
	item := &domain.Item{
		ID:          gofakeit.Number(1, 10000),
		UserID:      userID,
		Name:        gofakeit.Drink(),
		Description: gofakeit.SentenceSimple(),
		Price:       gofakeit.Number(5000, 10000),
		Cost:        gofakeit.Number(5000, 8000),
		Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
		Barcode:     gofakeit.Numerify("############"),
		ExpiryAt:    gofakeit.FutureDate(),
		Size:        domain.ItemSizeSmall,
		CreatedAt:   time.Now(),
	}
	err := item.Validate()
	assert.NoError(t, err)

	return item
}
