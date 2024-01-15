package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/text/language"

	"github.com/gin-gonic/gin"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/i18n"
	"github.com/psi59/payhere-assignment/internal/mocks/ucmocks"
	"github.com/psi59/payhere-assignment/usecase/authtoken"
	"github.com/psi59/payhere-assignment/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewAuthMiddleware(t *testing.T) {
	userUsecase := &user.Service{}
	authTokenUsecase := &authtoken.Service{}
	type args struct {
		userUsecase      user.Usecase
		authTokenUsecase authtoken.Usecase
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				userUsecase:      userUsecase,
				authTokenUsecase: authTokenUsecase,
			},
			wantErr: false,
		},
		{
			name: "nil userUsecase",
			args: args{
				userUsecase:      nil,
				authTokenUsecase: authTokenUsecase,
			},
			wantErr: true,
		},
		{
			name: "nil authTokenUsecase",
			args: args{
				userUsecase:      userUsecase,
				authTokenUsecase: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAuthMiddleware(tt.args.userUsecase, tt.args.authTokenUsecase)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestAuthMiddleware_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecase := ucmocks.NewMockUserUsecase(ctrl)
	authTokenUsecase := ucmocks.NewMockAuthTokenUsecase(ctrl)

	r := gin.New()
	authMiddleware, err := NewAuthMiddleware(userUsecase, authTokenUsecase)
	require.NoError(t, err)
	r.POST("/", authMiddleware.Auth(), func(ginCtx *gin.Context) {
		ginCtx.Status(http.StatusNoContent)
	})

	userDomain, err := domain.NewUser(
		gofakeit.Regex(`^01\d{8,9}$`),
		gofakeit.Password(true, true, true, true, true, 10),
		gofakeit.Date(),
	)
	require.NoError(t, err)
	userDomain.ID = gofakeit.Number(1, 10)
	token := gofakeit.UUID()
	httpRequest, err := http.NewRequest(http.MethodPost, "/", nil)
	require.NoError(t, err)
	httpRequest.Header.Set("Authorization", "Bearer "+token)

	t.Run("OK", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(&authtoken.VerifyOutput{
			Identifier: strconv.Itoa(userDomain.ID),
			ExpiresAt:  gofakeit.FutureDate(),
		}, nil)
		authTokenUsecase.EXPECT().GetBlacklist(gomock.Any(), &authtoken.GetBlacklistInput{
			Token: token,
		}).Return(nil, domain.ErrTokenBlacklistNotFound)

		userUsecase.EXPECT().Get(gomock.Any(), &user.GetInput{
			UserID: userDomain.ID,
		}).Return(&user.GetOutput{User: userDomain}, nil)

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		assert.Equal(t, http.StatusNoContent, responseWriter.Code)
	})

	t.Run("empty token", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		invalidHttpRequest, err := http.NewRequest(http.MethodPost, "/", nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, invalidHttpRequest)

		assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
	})

	t.Run("만료된 토큰", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(nil, domain.ErrExpiredToken)

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
		assert.Equal(t, http.StatusUnauthorized, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.ExpiredToken, nil), resp.Meta.Message)
	})

	t.Run("토큰 검증 시 알 수 없는 에러 발생", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(nil, gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
		assert.Equal(t, http.StatusUnauthorized, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.Unauthorized, nil), resp.Meta.Message)
	})

	t.Run("잘못된 토큰일 경우", func(t *testing.T) {
		expiresAt := gofakeit.FutureDate()
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(&authtoken.VerifyOutput{
			Identifier: gofakeit.UUID(),
			ExpiresAt:  expiresAt,
		}, nil)
		authTokenUsecase.EXPECT().GetBlacklist(gomock.Any(), &authtoken.GetBlacklistInput{
			Token: token,
		}).Return(&authtoken.GetBlacklistOutput{
			Token: &domain.AuthToken{
				Token:     token,
				ExpiresAt: expiresAt,
			},
		}, nil)

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})

	t.Run("토큰 블랙리스트 조회 실패", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(&authtoken.VerifyOutput{
			Identifier: strconv.Itoa(userDomain.ID),
			ExpiresAt:  gofakeit.FutureDate(),
		}, nil)
		authTokenUsecase.EXPECT().GetBlacklist(gomock.Any(), &authtoken.GetBlacklistInput{
			Token: token,
		}).Return(nil, gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})

	t.Run("회원이 존재하지 않을 때", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(&authtoken.VerifyOutput{
			Identifier: strconv.Itoa(userDomain.ID),
			ExpiresAt:  gofakeit.FutureDate(),
		}, nil)
		authTokenUsecase.EXPECT().GetBlacklist(gomock.Any(), &authtoken.GetBlacklistInput{
			Token: token,
		}).Return(nil, domain.ErrTokenBlacklistNotFound)

		userUsecase.EXPECT().Get(gomock.Any(), &user.GetInput{
			UserID: userDomain.ID,
		}).Return(nil, domain.ErrUserNotFound)

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
		assert.Equal(t, http.StatusUnauthorized, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.UserNotFound, nil), resp.Meta.Message)
	})

	t.Run("회원 조회 실패", func(t *testing.T) {
		authTokenUsecase.EXPECT().Verify(gomock.Any(), &authtoken.VerifyInput{
			Token: token,
		}).Return(&authtoken.VerifyOutput{
			Identifier: strconv.Itoa(userDomain.ID),
			ExpiresAt:  gofakeit.FutureDate(),
		}, nil)
		authTokenUsecase.EXPECT().GetBlacklist(gomock.Any(), &authtoken.GetBlacklistInput{
			Token: token,
		}).Return(nil, domain.ErrTokenBlacklistNotFound)

		userUsecase.EXPECT().Get(gomock.Any(), &user.GetInput{
			UserID: userDomain.ID,
		}).Return(nil, gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}
