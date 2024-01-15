package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/psi59/payhere-assignment/internal/ginhelper"

	"github.com/psi59/payhere-assignment/internal/i18n"
	"golang.org/x/text/language"

	"github.com/psi59/payhere-assignment/domain"
	"github.com/stretchr/testify/assert"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/gin-gonic/gin"

	"github.com/psi59/payhere-assignment/internal/mocks/ucmocks"

	"github.com/psi59/payhere-assignment/usecase/authtoken"
	"github.com/psi59/payhere-assignment/usecase/user"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewUserHandler(t *testing.T) {
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
			got, err := NewUserHandler(tt.args.userUsecase, tt.args.authTokenUsecase)
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

func TestUserHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecase := ucmocks.NewMockUserUsecase(ctrl)
	authTokenUsecase := ucmocks.NewMockAuthTokenUsecase(ctrl)
	userDomain := newTestUser(t, gofakeit.Password(true, true, true, true, true, 10))

	r := gin.New()
	handler, err := NewUserHandler(userUsecase, authTokenUsecase)
	require.NoError(t, err)
	r.POST("/", handler.SignUp)

	t.Run("OK", func(t *testing.T) {
		signUpRequest := SignUpRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    userDomain.Password,
		}
		userUsecase.EXPECT().Create(gomock.Any(), &user.CreateInput{
			PhoneNumber: signUpRequest.PhoneNumber,
			Password:    signUpRequest.Password,
		}).Return(&user.CreateOutput{User: userDomain}, nil)
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signUpRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		assert.Equal(t, http.StatusNoContent, responseWriter.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(map[string]any{
			"phoneNumber": 10,
		})
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("invalid request", func(t *testing.T) {
		signUpRequest := SignUpRequest{
			PhoneNumber: gofakeit.UUID(),
			Password:    userDomain.Password,
		}
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signUpRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("user already exists", func(t *testing.T) {
		signUpRequest := SignUpRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    userDomain.Password,
		}
		userUsecase.EXPECT().Create(gomock.Any(), &user.CreateInput{
			PhoneNumber: signUpRequest.PhoneNumber,
			Password:    signUpRequest.Password,
		}).Return(nil, domain.ErrUserAlreadyExists)
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signUpRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, responseWriter.Code)
		assert.Equal(t, http.StatusConflict, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.UserAlreadyExists, nil), resp.Meta.Message)
	})

	t.Run("unexpected usecase error", func(t *testing.T) {
		signUpRequest := SignUpRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    userDomain.Password,
		}
		userUsecase.EXPECT().Create(gomock.Any(), &user.CreateInput{
			PhoneNumber: signUpRequest.PhoneNumber,
			Password:    signUpRequest.Password,
		}).Return(nil, gofakeit.Error())
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signUpRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}

func TestUserHandler_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecase := ucmocks.NewMockUserUsecase(ctrl)
	authTokenUsecase := ucmocks.NewMockAuthTokenUsecase(ctrl)

	plainPassword := gofakeit.Password(true, true, true, true, true, 10)
	userDomain := newTestUser(t, plainPassword)
	userDomain.ID = gofakeit.Number(1, 10)

	r := gin.New()
	handler, err := NewUserHandler(userUsecase, authTokenUsecase)
	require.NoError(t, err)
	r.POST("/", handler.SignIn)

	t.Run("OK", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    plainPassword,
		}
		userUsecase.EXPECT().GetByPhoneNumber(gomock.Any(), &user.GetByPhoneNumberInput{
			PhoneNumber: signInRequest.PhoneNumber,
		}).Return(&user.GetOutput{User: userDomain}, nil)
		authTokenUsecase.EXPECT().Create(gomock.Any(), &authtoken.CreateInput{
			Identifier: strconv.Itoa(userDomain.ID),
		}).Return(&authtoken.CreateOutput{
			Token:     gofakeit.UUID(),
			ExpiresAt: gofakeit.FutureDate(),
		}, nil)

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		assert.Equal(t, http.StatusOK, responseWriter.Code)
	})

	t.Run("바인딩 에러", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(map[string]any{
			"password": 1,
		})
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("잘못된 요청", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    "",
		}

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("존재하지 않는 유저", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    plainPassword,
		}
		userUsecase.EXPECT().GetByPhoneNumber(gomock.Any(), &user.GetByPhoneNumberInput{
			PhoneNumber: signInRequest.PhoneNumber,
		}).Return(nil, domain.ErrUserNotFound)

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, responseWriter.Code)
		assert.Equal(t, http.StatusNotFound, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.UserNotFound, nil), resp.Meta.Message)
	})

	t.Run("유저 조회 실패", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    plainPassword,
		}
		userUsecase.EXPECT().GetByPhoneNumber(gomock.Any(), &user.GetByPhoneNumberInput{
			PhoneNumber: signInRequest.PhoneNumber,
		}).Return(nil, gofakeit.ErrorDatabase())

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})

	t.Run("비밀번호 불일치", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    gofakeit.UUID(),
		}
		userUsecase.EXPECT().GetByPhoneNumber(gomock.Any(), &user.GetByPhoneNumberInput{
			PhoneNumber: signInRequest.PhoneNumber,
		}).Return(&user.GetOutput{User: userDomain}, nil)

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.PasswordMismatch, nil), resp.Meta.Message)
	})

	t.Run("토큰 생성 실패", func(t *testing.T) {
		signInRequest := SignInRequest{
			PhoneNumber: userDomain.PhoneNumber,
			Password:    plainPassword,
		}
		userUsecase.EXPECT().GetByPhoneNumber(gomock.Any(), &user.GetByPhoneNumberInput{
			PhoneNumber: signInRequest.PhoneNumber,
		}).Return(&user.GetOutput{User: userDomain}, nil)
		authTokenUsecase.EXPECT().Create(gomock.Any(), &authtoken.CreateInput{
			Identifier: strconv.Itoa(userDomain.ID),
		}).Return(nil, gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(signInRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}

func TestUserHandler_SignOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecase := ucmocks.NewMockUserUsecase(ctrl)
	authTokenUsecase := ucmocks.NewMockAuthTokenUsecase(ctrl)

	r := gin.New()
	handler, err := NewUserHandler(userUsecase, authTokenUsecase)
	require.NoError(t, err)
	r.POST("/", handler.SignOut)

	t.Run("OK", func(t *testing.T) {
		token := gofakeit.UUID()
		authTokenUsecase.EXPECT().RegisterBlacklist(gomock.Any(), &authtoken.RegisterBlacklistInput{
			Token: token,
		}).Return(nil)

		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodPost, "/", nil)
		require.NoError(t, err)
		httpRequest.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(responseWriter, httpRequest)

		assert.Equal(t, http.StatusNoContent, responseWriter.Code)
	})

	t.Run("이미 등록된 토큰일 경우", func(t *testing.T) {
		token := gofakeit.UUID()
		authTokenUsecase.EXPECT().RegisterBlacklist(gomock.Any(), &authtoken.RegisterBlacklistInput{
			Token: token,
		}).Return(domain.ErrTokenBlacklistAlreadyExists)

		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodPost, "/", nil)
		require.NoError(t, err)
		httpRequest.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
		assert.Equal(t, http.StatusUnauthorized, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.TokenBlacklistAlreadyExists, nil), resp.Meta.Message)
	})

	t.Run("예상하지 못한 토큰 등록 에러", func(t *testing.T) {
		token := gofakeit.UUID()
		authTokenUsecase.EXPECT().RegisterBlacklist(gomock.Any(), &authtoken.RegisterBlacklistInput{
			Token: token,
		}).Return(gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodPost, "/", nil)
		require.NoError(t, err)
		httpRequest.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(responseWriter, httpRequest)

		var resp ginhelper.Response
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}

func newTestUser(t *testing.T, password string) *domain.User {
	userDomain, err := domain.NewUser(
		gofakeit.Regex(`^01\d{8,9}$`),
		password,
		gofakeit.Date(),
	)
	require.NoError(t, err)
	userDomain.ID = gofakeit.Number(1, 10)

	return userDomain
}
