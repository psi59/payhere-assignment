package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/psi59/gopkg/ctxlog"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/rs/zerolog"
)

func SetContext() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginhelper.SetContext(ginCtx, ginCtx)
		ginCtx.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginhelper.GetContext(ginCtx)
		req := ginCtx.Request
		ctx = ctxlog.WithLogger(ctx)
		ginhelper.SetContext(ginCtx, ctx)
		t := time.Now()

		var lf logFields
		var logLevel zerolog.Level
		var handlerErr error
		lf.RequestID = requestid.Get(ginCtx)
		lf.IP = ginCtx.ClientIP()
		lf.Method = req.Method
		lf.URI = req.RequestURI
		lf.ContentLength = req.ContentLength
		lf.RequestTime = t
		defer func(lvl *zerolog.Level, lf *logFields) {
			res := ginCtx.Writer
			lf.ResponseSize = res.Size()
			lf.StatusCode = res.Status()
			lf.Duration = time.Since(t).String()

			ll := ctxlog.GetLogger(ctx)
			tmp := ll.Logger()
			logger := tmp.WithLevel(*lvl).Interface("requestInfo", lf)
			if !valid.IsNil(handlerErr) {
				logger.Err(handlerErr).Str("stackTrace", fmt.Sprintf("%+v", handlerErr))
			}

			logger.Send()
		}(&logLevel, &lf)

		ginCtx.Next()
		if len(ginCtx.Errors) == 0 {
			logLevel = zerolog.InfoLevel
			return
		}

		handlerErr = ginCtx.Errors.Last().Err
		var httpError *domain.HTTPError
		if !errors.As(handlerErr, &httpError) {
			logLevel = zerolog.ErrorLevel
			ginCtx.JSON(http.StatusInternalServerError, domain.Response{
				Meta: domain.ResponseMeta{
					Code:    http.StatusInternalServerError,
					Message: "",
				},
			})
			return
		}

		logLevel = zerolog.WarnLevel
		ginCtx.JSON(
			httpError.StatusCode,
			domain.Response{
				Meta: domain.ResponseMeta{
					Code:    httpError.StatusCode,
					Message: httpError.Message(),
				},
			},
		)
	}
}

type logFields struct {
	RequestID     string    `json:"requestId"`
	IP            string    `json:"ip"`
	Method        string    `json:"method"`
	URI           string    `json:"uri"`
	ContentLength int64     `json:"contentLength"`
	RequestTime   time.Time `json:"requestTime"`
	StatusCode    int       `json:"statusCode"`
	Duration      string    `json:"duration"`
	ResponseSize  int       `json:"responseSize"`
}
