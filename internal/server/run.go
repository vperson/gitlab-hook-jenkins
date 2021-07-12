package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vperson/gitlab-hook-jenkins/internal/server/router"
	"net/http"
)

func Run(debug bool, host string, port int) error {
	var (
		runMode = ""
	)
	if debug == true {
		runMode = "debug"
	} else {
		runMode = "release"
	}
	gin.SetMode(runMode)

	routersInit := router.InitRouter()
	endpoint := fmt.Sprintf("%s:%d", host, port)
	maxHeaderBytes := 1 << 20

	s := &http.Server{
		Addr:           endpoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Infof("start http server listening %s", endpoint)

	return s.ListenAndServe()
}
