package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/psi59/payhere-assignment/usecase/authtoken"
	"github.com/psi59/payhere-assignment/usecase/user"
)

type Authenticator struct {
	userUsecase      user.Usecase
	authTokenUsecase authtoken.Usecase
}

func NewAuthenticator(userUsecase user.Usecase, authTokenUsecase authtoken.Usecase) (*Authenticator, error) {
	switch {
	case valid.IsNil(userUsecase):
		return nil, user.ErrNilUsecase
	case valid.IsNil(authTokenUsecase):
		return nil, authtoken.ErrNilUsecase
	}

	return &Authenticator{
		userUsecase:      userUsecase,
		authTokenUsecase: authTokenUsecase,
	}, nil
}

func (a *Authenticator) Auth() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginhelper.GetContext(ginCtx)
		token := ginhelper.GetToken(ginCtx)
		if len(token) == 0 {
			ginhelper.Error(ginCtx, errors.WithStack(ginhelper.NewHTTPError(http.StatusUnauthorized, i18n.Unauthorized, fmt.Errorf("empty token"))))
			ginCtx.Abort()
			return
		}

		verifyTokenOutput, err := a.authTokenUsecase.Verify(ctx, &authtoken.VerifyInput{
			Token: token,
		})
		if err != nil {
			msgID := i18n.Unauthorized
			if errors.Is(err, domain.ErrExpiredToken) {
				msgID = i18n.ExpiredToken
			}

			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusUnauthorized, msgID, errors.WithStack(err)))
			ginCtx.Abort()
			return
		}

		userID, err := strconv.Atoi(verifyTokenOutput.Identifier)
		if err != nil {
			ginhelper.Error(ginCtx, errors.Wrap(err, "failed to parse userID"))
			ginCtx.Abort()
			return
		}

		userGetOutput, err := a.userUsecase.Get(ctx, &user.GetInput{
			UserID: userID,
		})
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusUnauthorized, i18n.UserNotFound, errors.WithStack(err)))
				ginCtx.Abort()
				return
			}

			ginhelper.Error(ginCtx, errors.Wrap(err, "failed to get user"))
			ginCtx.Abort()
			return
		}
		ctx = context.WithValue(ctx, domain.CtxKeyUser, userGetOutput.User)
		ginhelper.SetContext(ginCtx, ctx)

		ginCtx.Next()
	}
}
