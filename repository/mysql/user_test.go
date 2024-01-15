package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jinzhu/copier"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	ctx := db.ContextWithConn(context.TODO(), conn)

	t.Run("OK", func(t *testing.T) {
		user := newTestUser(t)
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		require.True(t, user.ID > 0)
	})

	t.Run("nil Context", func(t *testing.T) {
		user := newTestUser(t)
		err := repo.Create(nil, user)
		require.Error(t, err)
	})

	t.Run("nil user", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		require.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		user := newTestUser(t)
		err := repo.Create(context.TODO(), user)
		require.Error(t, err)
	})

	t.Run("Duplicated PhoneNumber", func(t *testing.T) {
		user := newTestUser(t)
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		require.True(t, user.ID > 0)

		var dupl domain.User
		err = copier.Copy(&dupl, &user)
		dupl.ID = 0
		require.NoError(t, err)
		err = repo.Create(ctx, &dupl)
		require.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})
}

func newTestUser(t *testing.T) *domain.User {
	user, err := domain.NewUser(
		gofakeit.Regex(`^01\d{8,9}$`),
		gofakeit.Password(true, true, true, true, true, 72),
		time.Unix(time.Now().Unix(), 0).UTC(),
	)
	require.NoError(t, err)
	return user
}

func TestUserRepository_Get(t *testing.T) {
	repo := NewUserRepository()
	ctx := db.ContextWithConn(context.TODO(), conn)

	user := newTestUser(t)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		got, err := repo.Get(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, user, got)
	})

	t.Run("nil Context", func(t *testing.T) {
		got, err := repo.Get(nil, user.ID)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid userID", func(t *testing.T) {
		got, err := repo.Get(ctx, gofakeit.IntRange(-10, 0))
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := repo.Get(context.TODO(), user.ID)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		got, err := repo.Get(ctx, gofakeit.IntRange(10000, 20000))
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestUserRepository_GetByPhoneNumber(t *testing.T) {
	repo := NewUserRepository()
	ctx := db.ContextWithConn(context.TODO(), conn)

	user := newTestUser(t)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		got, err := repo.GetByPhoneNumber(ctx, user.PhoneNumber)
		require.NoError(t, err)
		require.Equal(t, user, got)
	})

	t.Run("nil Context", func(t *testing.T) {
		got, err := repo.GetByPhoneNumber(nil, user.PhoneNumber)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("empty phoneNumber", func(t *testing.T) {
		got, err := repo.GetByPhoneNumber(ctx, "")
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := repo.GetByPhoneNumber(context.TODO(), user.PhoneNumber)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		got, err := repo.GetByPhoneNumber(ctx, gofakeit.Phone())
		require.Error(t, err)
		require.Nil(t, got)
	})
}
