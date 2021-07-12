package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type RequestOption struct {
	ctx *gin.Context
}

func (r *RequestOption) ParseJSON(obj interface{}) error {
	if err := r.ctx.ShouldBindJSON(obj); err != nil {
		return fmt.Errorf("解析请求参数发生错误 - %s", err.Error())
	}
	return nil
}
