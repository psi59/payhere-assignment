package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/psi59/gopkg/ctxlog"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/i18n"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/psi59/payhere-assignment/usecase/authtoken"
	"github.com/psi59/payhere-assignment/usecase/user"
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

func (h *UserHandler) SignUp(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)
	var req SignUpRequest
	if err := ginCtx.BindJSON(&req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	if err := req.Validate(); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	userCreateOutput, err := h.userUsecase.Create(ctx, &user.CreateInput{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusConflict, i18n.UserAlreadyExists, errors.WithStack(err)))
		default:
			ginhelper.Error(ginCtx, errors.WithStack(err))
		}

		return
	}
	ctxlog.WithInt(ctx, "userID", userCreateOutput.User.ID)

	ginCtx.Status(http.StatusNoContent)
}

func (h *UserHandler) SignIn(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)
	var req SignInRequest
	if err := ginCtx.BindJSON(&req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	if err := valid.ValidateStruct(req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}

	userGetOutput, err := h.userUsecase.GetByPhoneNumber(ctx, &user.GetByPhoneNumberInput{PhoneNumber: req.PhoneNumber})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusNotFound, i18n.UserNotFound, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}
	userDomain := userGetOutput.User
	if err := userDomain.ComparePassword(req.Password); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.PasswordMismatch, errors.WithStack(err)))
		return
	}

	createTokenOutput, err := h.authTokenUsecase.Create(ctx, &authtoken.CreateInput{Identifier: strconv.Itoa(userDomain.ID)})
	if err != nil {
		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}

	ginhelper.Success(ginCtx, SignInResponse{
		Token:     createTokenOutput.Token,
		ExpiresAt: createTokenOutput.ExpiresAt,
	})
}

func (h *UserHandler) SignOut(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)
	token := ginhelper.GetToken(ginCtx)

	if err := h.authTokenUsecase.RegisterBlacklist(ctx, &authtoken.RegisterBlacklistInput{Token: token}); err != nil {
		if errors.Is(err, domain.ErrTokenBlacklistAlreadyExists) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusUnauthorized, i18n.TokenBlacklistAlreadyExists, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}

	ginCtx.Status(http.StatusNoContent)
	return
}

type SignUpRequest struct {
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
