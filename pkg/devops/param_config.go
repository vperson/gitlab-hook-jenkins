package devops

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var (
	projectsCount map[string]int
	projects      map[string]int
	config        *viper.Viper
	configDir     = "storage"
	configFile    = "data.yaml"
	BranchMapJob  BranchMapJobOptions
	mux           sync.RWMutex
)

type BranchMapJobOptions struct {
	BranchesMapJob []struct {
		// gitlab project
		Project string `json:"project" yaml:"project"`
		Config  []struct {
			// gitlab branch
			Branch string `json:"branch" yaml:"branch"`
			// jenkins build params
			Params []struct {
				Name  string `json:"name" yaml:"name"`
				Value string `json:"value" yaml:"value"`
			} `json:"params" yaml:"params"`
		} `json:"config" yaml:"config"`

		// jenkins job
		JenkinsJob string `json:"jenkinsJob" yaml:"jenkinsJob"`
	} `json:"branchesMapJob" yaml:"branchesMapJob"`
}

func ConfigInit() {
	config = viper.New()
	BranchMapJob = BranchMapJobOptions{}
	projectsCount = make(map[string]int)
	config.AddConfigPath(configDir)
	config.SetConfigName(configFile)
	config.SetConfigType("yaml")
}
func BranchMapJobConfigReload() error {
	mux.Lock()
	defer mux.Unlock()
	projectsCount = make(map[string]int)
	projects = make(map[string]int)

	err := config.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件错误 -- %v", err)
	}

	if err = config.Unmarshal(&BranchMapJob); err != nil {
		return fmt.Errorf("反序列化配置文件错误 -- %v", err)
	}

	for index, p := range BranchMapJob.BranchesMapJob {
		projectsCount[p.Project] += 1
		projects[p.Project] = index
	}

	checkProjectCount()
	fmt.Printf("xxx %+v\n", BranchMapJob.BranchesMapJob)
	return nil
}

// 对于出现多次的Projects打印警告日志
func checkProjectCount() {
	for k, v := range projectsCount {
		if v > 1 {
			log.Warnf("项目: %s 出现了 %d 次配置可能相互覆盖", k, v)
		}
	}
}

func ReadBranchesMapJob() (BranchMapJobOptions, map[string]int) {
	mux.RLock()
	defer mux.RUnlock()
	return BranchMapJob, projects
}

// 当文件发生变化的时候动态加载
func WatchBranchMapJobConfig() {
	config.WatchConfig()
	config.OnConfigChange(func(in fsnotify.Event) {
		log.Infof("config file changed: %s", in.Name)
		switch in.Op {
		case fsnotify.Write:
			err := BranchMapJobConfigReload()
			if err != nil {
				log.Errorf("重载配置文件错误 -- %v", err)
			}
			fmt.Println(projects)
			fmt.Println(BranchMapJob)
		}
	})
}
