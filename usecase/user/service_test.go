package user

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/mocks/repomocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	ctx := context.TODO()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			user.ID = gofakeit.Number(1, 100)
			return nil
		})
		input := &CreateInput{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Regex(`^01\d{8,9}$`),
			Password:    gofakeit.Password(true, true, true, true, true, 10),
		}
		got, err := srv.Create(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotNil(t, got.User)
	})

	t.Run("nil context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		input := &CreateInput{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Regex(`^01\d{8,9}$`),
			Password:    gofakeit.Password(true, true, true, true, true, 10),
		}
		got, err := srv.Create(nil, input)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		input := &CreateInput{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Regex(`^01\d{8,9}$`),
			Password:    gofakeit.Password(true, true, false, true, true, 10),
		}
		got, err := srv.Create(ctx, input)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		input := &CreateInput{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Regex(`^01\d{8,9}$`),
			Password:    gofakeit.Password(true, true, false, true, true, 10),
		}
		got, err := srv.Create(ctx, input)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("failed to create user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
			return gofakeit.Error()
		})
		input := &CreateInput{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Regex(`^01\d{8,9}$`),
			Password:    gofakeit.Password(true, true, true, true, true, 10),
		}
		got, err := srv.Create(ctx, input)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		user.ID = gofakeit.Number(1, 100)
		userRepo.EXPECT().Get(ctx, user.ID).DoAndReturn(func(ctx context.Context, i int) (*domain.User, error) {
			return user, nil
		})
		got, err := srv.Get(ctx, &GetInput{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, user, got.User)
	})

	t.Run("nil Context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		user.ID = gofakeit.Number(1, 100)
		got, err := srv.Get(nil, &GetInput{
			UserID: user.ID,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		got, err := srv.Get(ctx, nil)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		got, err := srv.Get(ctx, &GetInput{
			UserID: 0,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		user.ID = gofakeit.Number(1, 100)
		userRepo.EXPECT().Get(ctx, user.ID).DoAndReturn(func(ctx context.Context, i int) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		})
		got, err := srv.Get(ctx, &GetInput{UserID: user.ID})
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestService_GetByPhoneNumber(t *testing.T) {
	ctx := context.TODO()

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		userRepo.EXPECT().GetByPhoneNumber(ctx, user.PhoneNumber).DoAndReturn(func(ctx context.Context, s string) (*domain.User, error) {
			return user, nil
		})
		got, err := srv.GetByPhoneNumber(ctx, &GetByPhoneNumberInput{
			PhoneNumber: user.PhoneNumber,
		})
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, user, got.User)
	})

	t.Run("nil Context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		got, err := srv.GetByPhoneNumber(nil, &GetByPhoneNumberInput{
			PhoneNumber: user.PhoneNumber,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		got, err := srv.GetByPhoneNumber(ctx, nil)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		got, err := srv.GetByPhoneNumber(ctx, &GetByPhoneNumberInput{
			PhoneNumber: "",
		})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repomocks.NewMockUserRepository(ctrl)
		srv, err := NewService(userRepo)
		require.NoError(t, err)

		user, err := domain.NewUser(
			gofakeit.Name(),
			gofakeit.Regex(`^01\d{8,9}$`),
			gofakeit.Password(true, true, true, true, true, 10),
			time.Unix(time.Now().Unix(), 0).UTC(),
		)
		require.NoError(t, err)
		userRepo.EXPECT().GetByPhoneNumber(ctx, user.PhoneNumber).DoAndReturn(func(ctx context.Context, s string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		})
		got, err := srv.GetByPhoneNumber(ctx, &GetByPhoneNumberInput{
			PhoneNumber: user.PhoneNumber,
		})
		require.Error(t, err)
		require.Nil(t, got)
	})
}
