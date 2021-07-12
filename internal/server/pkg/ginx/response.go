package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ResponseOption struct {
	ctx *gin.Context
}

func (r *ResponseOption) Success() {
	r.ctx.JSON(
		http.StatusOK,
		"ok",
	)
}

func (r *ResponseOption) StatusInternalServerError(err error) {
	log.Error(err)
	r.ctx.JSON(
		http.StatusInternalServerError,
		"fail",
	)
}

func (r *ResponseOption) ParameterError(err error) {
	r.ctx.JSON(
		http.StatusBadRequest,
		fmt.Sprintf("parameter error: %v", err),
	)
}
