package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/psi59/payhere-assignment/repository"

	"github.com/jinzhu/copier"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"

	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestItemRepository_Create(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		assert.True(t, item.ID > 0)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(nil, item)
		assert.Error(t, err)
	})

	t.Run("nil item", func(t *testing.T) {
		err := itemRepo.Create(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(context.TODO(), item)
		assert.Error(t, err)
	})

	t.Run("invalid item", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		item.UserID = 0
		err := itemRepo.Create(ctx, item)
		assert.Error(t, err)
	})

	t.Run("중복 아이템", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		var dupl domain.Item
		err := copier.Copy(&dupl, item)
		assert.NoError(t, err)

		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		assert.True(t, item.ID > 0)

		err = itemRepo.Create(ctx, &dupl)
		assert.ErrorIs(t, err, domain.ErrItemAlreadyExists)
	})
}

func TestItemRepository_Get(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()
	item := newTestItem(t, user.ID)
	err = itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)
		assert.Equal(t, item, got)
	})

	t.Run("nil context", func(t *testing.T) {
		got, err := itemRepo.Get(nil, item.UserID, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid userID", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, 0, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid itemID", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, item.UserID, 0)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := itemRepo.Get(context.TODO(), item.UserID, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("item not found", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, gofakeit.Number(1000, 2000), item.ID)
		assert.Error(t, err, domain.ErrItemNotFound)
		assert.Nil(t, got)

		got, err = itemRepo.Get(ctx, item.UserID, gofakeit.Number(1000, 2000))
		assert.Error(t, err, domain.ErrItemNotFound)
		assert.Nil(t, got)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()
	item := newTestItem(t, user.ID)
	err = itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		err := itemRepo.Delete(ctx, item.UserID, item.ID)
		assert.NoError(t, err)
	})

	t.Run("nil context", func(t *testing.T) {
		err := itemRepo.Delete(nil, item.UserID, item.ID)
		assert.Error(t, err)

	})

	t.Run("invalid userID", func(t *testing.T) {
		err := itemRepo.Delete(ctx, 0, item.ID)
		assert.Error(t, err)

	})

	t.Run("invalid itemID", func(t *testing.T) {
		err := itemRepo.Delete(ctx, item.UserID, 0)
		assert.Error(t, err)

	})

	t.Run("context without conn", func(t *testing.T) {
		err := itemRepo.Delete(context.TODO(), item.UserID, item.ID)
		assert.Error(t, err)

	})
}

func TestItemRepository_Update(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		description := gofakeit.SentenceSimple()
		price := gofakeit.Number(1000, 10000)
		cost := gofakeit.Number(1000, 10000)
		category := gofakeit.SentenceSimple()
		barcode := gofakeit.RandomString([]string{"coffee", "tea", "desert"})
		size := domain.ItemSizeLarge
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:        &name,
			Description: &description,
			Price:       &price,
			Cost:        &cost,
			Category:    &category,
			Barcode:     &barcode,
			Size:        &size,
			ExpiryAt:    &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.NoError(t, err)

		var expected domain.Item
		err = copier.Copy(&expected, item)
		assert.NoError(t, err)
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)

		expected.Name = name
		expected.Description = description
		expected.Price = price
		expected.Cost = cost
		expected.Category = category
		expected.Barcode = barcode
		expected.Size = size
		expected.ExpiryAt = expiryAt

		assert.Equal(t, &expected, got)
	})

	t.Run("부분 업데이트", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.NoError(t, err)

		var expected domain.Item
		err = copier.Copy(&expected, item)
		assert.NoError(t, err)
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)

		expected.Name = name
		expected.ExpiryAt = expiryAt

		assert.Equal(t, &expected, got)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}

		err = itemRepo.Update(nil, item.UserID, item.ID, input)
		assert.Error(t, err)

	})

	t.Run("invalid userID", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, 0, item.ID, input)
		assert.Error(t, err)
	})

	t.Run("invalid itemID", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, 0, input)
		assert.Error(t, err)
	})

	t.Run("nil input", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		err := itemRepo.Update(ctx, item.UserID, item.ID, nil)
		assert.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		err := itemRepo.Update(ctx, item.UserID, item.ID, &repository.UpdateItemInput{})
		assert.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(context.TODO(), item.UserID, item.ID, input)
		assert.Error(t, err)
	})

	t.Run("이름 중복", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		item2 := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item2)
		assert.NoError(t, err)

		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &item2.Name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrItemAlreadyExists)
	})
}

func newTestItem(t *testing.T, userID int) *domain.Item {
	item, err := domain.NewItem(
		userID,
		gofakeit.UUID(),
		gofakeit.SentenceSimple(),
		gofakeit.Number(5000, 10000),
		gofakeit.Number(3000, 5000),
		gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
		gofakeit.Numerify("##################"),
		time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		domain.ItemSize(gofakeit.RandomString([]string{string(domain.ItemSizeSmall), string(domain.ItemSizeLarge)})),
	)
	assert.NoError(t, err)
	item.CreatedAt = time.Unix(time.Now().Unix(), 0).UTC()

	return item
}
