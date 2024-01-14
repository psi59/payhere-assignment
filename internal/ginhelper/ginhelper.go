package ginhelper

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

const ctxKey = "_ginhelper/ctxKey"

func GetContext(c *gin.Context) context.Context {
	rawCtx, exists := c.Get(ctxKey)
	if exists {
		return rawCtx.(context.Context)
	}
	ctx := context.Context(c)
	SetContext(c, ctx)

	return ctx
}

func SetContext(ginCtx *gin.Context, ctx context.Context) {
	ginCtx.Set(ctxKey, ctx)
}

func GetToken(ginCtx *gin.Context) string {
	return strings.TrimPrefix(ginCtx.GetHeader("Authorization"), "Bearer ")
}
