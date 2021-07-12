package ginx

import "github.com/gin-gonic/gin"

type Option struct {
	ctx      *gin.Context
	Request  *RequestOption
	Response *ResponseOption
}

func New(ctx *gin.Context) *Option {
	return &Option{
		ctx: ctx,
		Request: &RequestOption{
			ctx: ctx,
		},
		Response: &ResponseOption{
			ctx: ctx,
		},
	}
}
