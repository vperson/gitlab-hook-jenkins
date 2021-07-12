package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vperson/gitlab-hook-jenkins/internal/server"
	"github.com/vperson/gitlab-hook-jenkins/pkg/bootstrap"
)

var (
	appArgs *bootstrap.AppArgs

	rootCmd = &cobra.Command{
		Use:          "hook",
		Short:        "gitlab hook jenkins",
		Long:         "gitlab system hook to post jenkins build",
		SilenceUsage: true,
	}

	appCmd = &cobra.Command{
		Use:   "server",
		Short: "start http server monitor hook for gitlab",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// 读取gitlab 项目 分支与jenkins job的配置
			err := bootstrap.ReadConfigData()
			if err != nil {
				return err
			}

			// 初始化Jenkins客户端
			err = bootstrap.NewJenkins(ctx, appArgs.Jenkins.URL, appArgs.Jenkins.Username, appArgs.Jenkins.Password)
			if err != nil {
				return err
			}
			return server.Run(appArgs.Debug, appArgs.Server.Host, appArgs.Server.Port)
		},
	}
)

func init() {
	appArgs = bootstrap.NewAppArgs()

	rootCmd.PersistentFlags().BoolVarP(&appArgs.Debug, "debug", "d", false, "enable debug mode")

	appCmd.PersistentFlags().StringVar(&appArgs.Jenkins.URL, "jenkins.url", "", "jenkins server url")
	appCmd.PersistentFlags().StringVar(&appArgs.Jenkins.Username, "jenkins.username", "", "jenkins login username")
	appCmd.PersistentFlags().StringVar(&appArgs.Jenkins.Password, "jenkins.password", "", "jenkins login password")

	appCmd.PersistentFlags().StringVar(&appArgs.Server.Host, "server.host", "0.0.0.0", "http listening address")
	appCmd.PersistentFlags().IntVar(&appArgs.Server.Port, "server.port", 8080, "http listening port")

	rootCmd.AddCommand(appCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
