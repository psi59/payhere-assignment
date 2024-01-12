package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/psi59/payhere-assignment/internal/i18n"

	"github.com/psi59/payhere-assignment/usecase/authtoken"

	"github.com/psi59/payhere-assignment/internal/ginhelper"

	"github.com/psi59/payhere-assignment/internal/ctxlog"

	"github.com/psi59/payhere-assignment/usecase/user"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/gin-gonic/gin"
	"github.com/psi59/payhere-assignment/domain"
)

type UserHandler struct {
	userUsecase      user.Usecase
	authTokenUsecase authtoken.Usecase
}

func NewUserHandler(userUsecase user.Usecase, authTokenUsecase authtoken.Usecase) (*UserHandler, error) {
	if valid.IsNil(userUsecase) {
		return nil, user.ErrNilUsecase
	}
	if valid.IsNil(authTokenUsecase) {
		return nil, authtoken.ErrNilUsecase
	}

	return &UserHandler{
		userUsecase:      userUsecase,
		authTokenUsecase: authTokenUsecase,
	}, nil
}

func (h *UserHandler) SignUp(c *gin.Context) {
	ctx := ginhelper.GetContext(c)
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, err)))
		return
	}
	if err := req.Validate(); err != nil {
		_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, err)))
		return
	}
	userDomain, err := h.userUsecase.Create(ctx, &user.CreateInput{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicatedUser):
			_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusConflict, i18n.UserAlreadyExists, err)))
		default:
			_ = c.Error(errors.WithStack(err))
		}

		return
	}
	ctxlog.WithInt(ctx, "userID", userDomain.ID)

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	ctx := ginhelper.GetContext(c)
	var req SignInRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, err)))
		return
	}
	if err := valid.ValidateStruct(req); err != nil {
		_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, err)))
		return
	}

	userDomain, err := h.userUsecase.GetByPhoneNumber(ctx, &user.GetByPhoneNumberInput{PhoneNumber: req.PhoneNumber})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusNotFound, i18n.UserNotFound, err)))
			return
		}

		_ = c.Error(errors.WithStack(err))
		return
	}
	if err := userDomain.ComparePassword(req.Password); err != nil {
		_ = c.Error(errors.WithStack(domain.NewHTTPError(http.StatusBadRequest, i18n.PasswordMismatch, err)))
		return
	}

	createTokenOutput, err := h.authTokenUsecase.Create(ctx, &authtoken.CreateInput{Identifier: strconv.Itoa(userDomain.ID)})
	if err != nil {
		_ = c.Error(errors.WithStack(err))
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Meta: domain.ResponseMeta{
			Code:    200,
			Message: "ok",
		},
		Data: SignInResponse{
			Token:     createTokenOutput.Token,
			ExpiresAt: createTokenOutput.ExpiresAt,
		},
	})
}

func (h *UserHandler) SignOut(c *gin.Context) {
	ctx := ginhelper.GetContext(c)
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if len(token) > 0 {
		if err := h.authTokenUsecase.RegisterToBlacklist(ctx, &authtoken.RegisterToBlacklistInput{Token: token}); err != nil {
			// 토큰 블랙리스트 중복 등록일 경우 200 반환
			if !errors.Is(err, domain.ErrDuplicatedTokenBlacklist) {
				_ = c.Error(err)
				return
			}

			log.Warn().Err(err).Send()
		}
	}

	c.JSON(http.StatusOK, domain.Response{
		Meta: domain.ResponseMeta{
			Code:    http.StatusOK,
			Message: "ok",
		},
	})
	return
}

type SignUpRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func (r *SignUpRequest) Validate() error {
	if err := valid.ValidatePhoneNumber(r.PhoneNumber); err != nil {
		return fmt.Errorf("%w: %q", err, r.PhoneNumber)
	}

	if err := valid.ValidatePassword(r.Password); err != nil {
		return fmt.Errorf("%w: %q", err, r.Password)
	}

	return nil
}

type SignInRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type SignInResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
