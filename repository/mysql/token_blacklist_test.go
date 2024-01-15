package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
)

func TestTokenBlacklistRepository_Create(t *testing.T) {
	repo := NewTokenBlacklistRepository()
	ctx := db.ContextWithConn(context.TODO(), conn)

	t.Run("OK", func(t *testing.T) {
		token := newTestTokenBlacklist()
		err := repo.Create(ctx, token)
		require.NoError(t, err)
	})

	t.Run("nil Context", func(t *testing.T) {
		token := newTestTokenBlacklist()
		err := repo.Create(nil, token)
		require.Error(t, err)
	})

	t.Run("nil token", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		require.Error(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		err := repo.Create(ctx, &domain.AuthToken{})
		require.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		token := newTestTokenBlacklist()
		err := repo.Create(context.TODO(), token)
		require.Error(t, err)
	})

	t.Run("Duplicate Token", func(t *testing.T) {
		token := newTestTokenBlacklist()
		err := repo.Create(ctx, token)
		require.NoError(t, err)

		err = repo.Create(ctx, token)
		require.ErrorIs(t, err, domain.ErrTokenBlacklistAlreadyExists)
	})
}

func TestTokenBlacklistRepository_Get(t *testing.T) {
	repo := NewTokenBlacklistRepository()
	ctx := db.ContextWithConn(context.TODO(), conn)

	token := newTestTokenBlacklist()
	err := repo.Create(ctx, token)
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		got, err := repo.Get(ctx, token.Token)
		require.NoError(t, err)
		require.Equal(t, token, got)
	})

	t.Run("token not exists", func(t *testing.T) {
		got, err := repo.Get(ctx, gofakeit.UUID())
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("nil Context", func(t *testing.T) {
		got, err := repo.Get(nil, gofakeit.UUID())
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("empty token", func(t *testing.T) {
		got, err := repo.Get(ctx, "")
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := repo.Get(context.TODO(), gofakeit.UUID())
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func newTestTokenBlacklist() *domain.AuthToken {
	return &domain.AuthToken{
		Token:     gofakeit.UUID(),
		ExpiresAt: time.Unix(time.Now().Unix(), 0).UTC(),
	}
}
