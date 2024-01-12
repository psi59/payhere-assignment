package ginhelper

import (
	"context"

	"github.com/gin-gonic/gin"
)

const ctxKey = "ctxKey"

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
