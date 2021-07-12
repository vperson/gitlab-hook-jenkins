package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vperson/gitlab-hook-jenkins/internal/server/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/ping", api.Ping)
	apiv1 := r.Group("/api")
	{
		apiv1.POST("/gitlab/hook", api.GitlabHook)
	}

	return r
}
