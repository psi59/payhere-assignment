package authtoken

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/psi59/payhere-assignment/repository/mysql"

	"github.com/psi59/payhere-assignment/domain"

	"github.com/brianvoe/gofakeit/v6"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/psi59/payhere-assignment/internal/mocks/repomocks"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		repo := mysql.NewTokenBlacklistRepository()
		got, err := NewService(gofakeit.LetterN(10), repo)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("empty secret", func(t *testing.T) {
		repo := mysql.NewTokenBlacklistRepository()
		got, err := NewService("", repo)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("nil tokenBlacklistRepository", func(t *testing.T) {
		got, err := NewService(gofakeit.LetterN(10), nil)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestService_Create(t *testing.T) {
	ctx := context.TODO()
	secret := gofakeit.LetterN(100)
	id := gofakeit.UUID()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Create(ctx, &CreateInput{Identifier: id})
		require.NoError(t, err)
		require.NotEmpty(t, got)
	})

	t.Run("nil context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Create(nil, &CreateInput{Identifier: id})
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Create(ctx, nil)
		require.Error(t, err)
		require.Empty(t, got)

		got, err = srv.Create(ctx, &CreateInput{Identifier: ""})
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func TestService_Verify(t *testing.T) {
	ctx := context.TODO()
	secret := gofakeit.LetterN(100)
	id := gofakeit.UUID()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		createOutput, err := srv.Create(ctx, &CreateInput{Identifier: id})
		require.NoError(t, err)
		require.NotEmpty(t, createOutput)

		got, err := srv.Verify(ctx, &VerifyInput{Token: createOutput.Token})
		require.NoError(t, err)
		require.Equal(t, id, got.Identifier)
	})

	t.Run("nil context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Verify(nil, &VerifyInput{Token: gofakeit.LetterN(500)})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Verify(ctx, nil)
		require.Error(t, err)
		require.Nil(t, got)

		got, err = srv.Verify(ctx, &VerifyInput{Token: ""})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		got, err := srv.Verify(ctx, &VerifyInput{Token: gofakeit.Sentence(10)})
		require.Error(t, err)
		require.Nil(t, got)

		invalidToken := base64.StdEncoding.EncodeToString([]byte(gofakeit.LetterN(100)))
		got, err = srv.Verify(ctx, &VerifyInput{Token: invalidToken})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid secret", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		now := time.Unix(time.Now().Unix(), 0).UTC()
		claims := &jwt.RegisteredClaims{
			ID:        xid.New().String(),
			Issuer:    "payhere-assignment",
			Subject:   gofakeit.UUID(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 0, 1)),
		}
		token, err := srv.createJWT(claims, []byte(gofakeit.UUID()))
		require.NoError(t, err)
		got, err := srv.Verify(ctx, &VerifyInput{Token: token})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("expired token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		now := time.Unix(time.Now().Unix(), 0).UTC()
		claims := &jwt.RegisteredClaims{
			ID:        xid.New().String(),
			Issuer:    "payhere-assignment",
			Subject:   gofakeit.UUID(),
			IssuedAt:  jwt.NewNumericDate(now.AddDate(0, 0, -2)),
			ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 0, -1)),
		}
		token, err := srv.createJWT(claims, srv.secret)
		require.NoError(t, err)
		got, err := srv.Verify(ctx, &VerifyInput{Token: token})
		require.ErrorIs(t, err, domain.ErrExpiredToken)
		require.Nil(t, got)
	})
}

func TestService_RegisterBlacklist(t *testing.T) {
	ctx := context.TODO()
	secret := gofakeit.LetterN(100)
	id := gofakeit.UUID()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		createOutput, err := srv.Create(ctx, &CreateInput{Identifier: id})
		require.NoError(t, err)
		require.NotEmpty(t, createOutput)

		tokenBlacklistRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		err = srv.RegisterBlacklist(ctx, &RegisterBlacklistInput{Token: createOutput.Token})
		require.NoError(t, err)
	})

	t.Run("nil context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		err = srv.RegisterBlacklist(nil, &RegisterBlacklistInput{Token: gofakeit.UUID()})
		require.ErrorIs(t, err, domain.ErrNilContext)
	})

	t.Run("nil input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		err = srv.RegisterBlacklist(ctx, nil)
		require.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		err = srv.RegisterBlacklist(ctx, &RegisterBlacklistInput{Token: ""})
		require.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		now := time.Unix(time.Now().Unix(), 0).UTC()
		claims := &jwt.RegisteredClaims{
			ID:        xid.New().String(),
			Issuer:    "payhere-assignment",
			Subject:   gofakeit.UUID(),
			IssuedAt:  jwt.NewNumericDate(now.AddDate(0, 0, -2)),
			ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 0, -1)),
		}
		token, err := srv.createJWT(claims, srv.secret)
		require.NoError(t, err)
		err = srv.RegisterBlacklist(ctx, &RegisterBlacklistInput{Token: token})
		require.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		err = srv.RegisterBlacklist(ctx, &RegisterBlacklistInput{Token: gofakeit.UUID()})
		require.Error(t, err)
	})

	t.Run("failed to create record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
		srv, err := NewService(secret, tokenBlacklistRepo)
		require.NoError(t, err)

		createOutput, err := srv.Create(ctx, &CreateInput{Identifier: id})
		require.NoError(t, err)
		require.NotEmpty(t, createOutput)

		tokenBlacklistRepo.EXPECT().Create(ctx, gomock.Any()).Return(gofakeit.Error())
		err = srv.RegisterBlacklist(ctx, &RegisterBlacklistInput{Token: createOutput.Token})
		require.Error(t, err)
	})
}

func TestService_GetBlacklist(t *testing.T) {
	ctx := context.TODO()
	secret := gofakeit.LetterN(100)
	token := &domain.AuthToken{
		Token:     gofakeit.UUID(),
		ExpiresAt: gofakeit.FutureDate(),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenBlacklistRepo := repomocks.NewMockTokenBlacklistRepository(ctrl)
	srv, err := NewService(secret, tokenBlacklistRepo)
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		tokenBlacklistRepo.EXPECT().Get(ctx, token.Token).DoAndReturn(func(ctx context.Context, s string) (*domain.AuthToken, error) {
			return token, nil
		})

		got, err := srv.GetBlacklist(ctx, &GetBlacklistInput{
			Token: token.Token,
		})
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, token, got.Token)
	})

	t.Run("token not found", func(t *testing.T) {
		tokenBlacklistRepo.EXPECT().Get(ctx, token.Token).DoAndReturn(func(ctx context.Context, s string) (*domain.AuthToken, error) {
			return nil, domain.ErrTokenBlacklistNotFound
		})

		got, err := srv.GetBlacklist(ctx, &GetBlacklistInput{
			Token: token.Token,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("token not found", func(t *testing.T) {
		got, err := srv.GetBlacklist(nil, &GetBlacklistInput{
			Token: token.Token,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		got, err := srv.GetBlacklist(ctx, nil)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		got, err := srv.GetBlacklist(ctx, &GetBlacklistInput{
			Token: "",
		})
		require.Error(t, err)
		require.Nil(t, got)
	})
}
