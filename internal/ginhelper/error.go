package ginhelper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/internal/i18n"
	"golang.org/x/text/language"
)

func Error(ginCtx *gin.Context, err error) {
	_ = ginCtx.Error(err)
	var httpError *HTTPError
	if !errors.As(err, &httpError) {
		ginCtx.JSON(http.StatusInternalServerError, Response{
			Meta: ResponseMeta{
				Code:    http.StatusInternalServerError,
				Message: i18n.T(language.English, i18n.InternalError, nil),
			},
		})
		return
	}

	ginCtx.JSON(
		httpError.StatusCode,
		Response{
			Meta: ResponseMeta{
				Code:    httpError.StatusCode,
				Message: httpError.Message(),
			},
		},
	)
}

func Success(ginCtx *gin.Context, data any) {
	ginCtx.JSON(http.StatusOK, Response{
		Meta: ResponseMeta{
			Code:    200,
			Message: "ok",
		},
		Data: data,
	})
}

type HTTPError struct {
	StatusCode int
	ErrorCode  string
	Internal   error
}

func NewHTTPError(statusCode int, msgID string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		ErrorCode:  msgID,
		Internal:   err,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%s:%d] %s: %v", e.ErrorCode, e.StatusCode, e.Message(), e.Internal)
}

func (e *HTTPError) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		_, _ = fmt.Fprintf(s, "[%s:%d] %s: %+v", e.ErrorCode, e.StatusCode, e.Message(), e.Internal)
		return
	}
}

func (e *HTTPError) Message() string {
	return i18n.T(language.English, e.ErrorCode, nil)
}
