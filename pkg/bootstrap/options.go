package bootstrap

import (
	"context"
	"github.com/vperson/gitlab-hook-jenkins/pkg/devops"
)

// 服务启动配置
type AppArgs struct {
	Debug   bool `mapstructure:"debug" json:"debug"`
	Jenkins struct {
		URL      string `mapstructure:"url" json:"url"`
		Username string `mapstructure:"username" json:"username"`
		Password string `mapstructure:"password" json:"password"`
	}
	Server struct {
		Port int    `mapstructure:"port" json:"port"`
		Host string `mapstructure:"host" json:"host"`
	}
}

func NewAppArgs(initFuncs ...func(args *AppArgs)) *AppArgs {
	args := &AppArgs{}

	for _, fn := range initFuncs {
		fn(args)
	}

	return args
}

var Jenkins *devops.JenkinsOptions

func NewJenkins(ctx context.Context, url, username, password string) error {
	var err error
	Jenkins, err = devops.NewJenkins(ctx, url, username, password)
	return err
}

func ReadConfigData() error {
	devops.ConfigInit()
	defer devops.WatchBranchMapJobConfig()
	return devops.BranchMapJobConfigReload()
}
