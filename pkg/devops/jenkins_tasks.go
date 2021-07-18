package devops

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	JenkinsBuildIds chan int64
	JenkinsQueue    []int64
)

func JenkinsBuildStatus() {
	timer := time.NewTimer(time.Second * 60)
	JenkinsBuildIds = make(chan int64, 4096)
	JenkinsQueue = []int64{}
	for {
		select {
		case id, ok := <-JenkinsBuildIds:
			if ok != true {
				continue
			}
			JenkinsQueue = append(JenkinsQueue, id)
		case <-timer.C:
			for _, i := range JenkinsQueue {
				r, err := Jenkins.BuildStatus(i)
				if err != nil {
					log.Errorf("jenkins id %d query fail", i)
					continue
				}

				switch r {
				case "SUCCESS":
					fmt.Println("执行成功")
				case "RUNNING":
					fmt.Println("运行中")
				default:
					fmt.Println("失败")
				}
			}
		}
	}
}
