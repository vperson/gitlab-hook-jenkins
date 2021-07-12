package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vperson/gitlab-hook-jenkins/internal/server/controllers"
	"github.com/vperson/gitlab-hook-jenkins/internal/server/pkg/ginx"
)

// 接受gitlab的系统hook
func GitlabHook(c *gin.Context) {
	ctx := ginx.New(c)

	hook := controllers.GitlabSystemHook{}

	if err := ctx.Request.ParseJSON(&hook); err != nil {
		ctx.Response.StatusInternalServerError(err)
		return
	}
	if ok := hook.CheckEvent(); ok != true {
		ctx.Response.Success()
		return
	}
	if err := hook.PostJenkins(); err != nil {
		ctx.Response.ParameterError(err)
		return
	}

	ctx.Response.Success()
}

func Ping(c *gin.Context) {
	c.JSON(200, "pong")
}
