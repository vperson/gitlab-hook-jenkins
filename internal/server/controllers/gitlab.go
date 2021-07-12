package controllers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vperson/gitlab-hook-jenkins/pkg/bootstrap"
	"github.com/vperson/gitlab-hook-jenkins/pkg/devops"
	"time"
)

type GitlabSystemHook struct {
	EventName   string `json:"event_name"`
	Before      string `json:"before"`
	After       string `json:"after"`
	Ref         string `json:"ref"`
	CheckoutSha string `json:"checkout_sha"`
	UserID      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	UserEmail   string `json:"user_email"`
	UserAvatar  string `json:"user_avatar"`
	ProjectID   int    `json:"project_id"`
	Project     struct {
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHTTPURL      string `json:"git_http_url"`
		GitSSHURL       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
}

func (h *GitlabSystemHook) CheckEvent() bool {
	if h.EventName == "push" {
		return true
	}

	return false
}

func (h *GitlabSystemHook) PostJenkins() (err error) {
	// 判断Jenkins Client是否为空
	jk := bootstrap.Jenkins
	if jk == nil {
		return fmt.Errorf("jenkins client is nil")
	}

	// 判断事件是否为push_event
	if h.EventName != "push" {
		log.WithFields(log.Fields{
			"event":    h.EventName,
			"project":  h.Project.Name,
			"userName": h.UserName,
		}).Info("不支持的事件类型")
		return fmt.Errorf("unsupported event type")
	}

	branches, projects := devops.ReadBranchesMapJob()
	index, ok := projects[h.Project.Name]
	if ok != true {
		log.WithFields(log.Fields{
			"event":    h.EventName,
			"project":  h.Project.Name,
			"userName": h.UserName,
		}).Info("不支持的project")
		return fmt.Errorf("project %s is unsupported ", h.Project.Name)
	}

	value := branches.BranchesMapJob[index]
	job := value.JenkinsJob
	branch := ""
	params := make(map[string]string)
	for _, i := range value.Config {
		if i.Branch != h.Ref {
			continue
		}
		branch = h.Ref
		for _, j := range i.Params {
			params[j.Name] = j.Value
		}
	}

	if job == "" {
		log.WithFields(log.Fields{
			"event":    h.EventName,
			"project":  h.Project.Name,
			"userName": h.UserName,
		}).Info("jenkins job 不为能为空")
	}
	if branch == "" {
		log.WithFields(log.Fields{
			"event":    h.EventName,
			"project":  h.Project.Name,
			"userName": h.UserName,
		}).Infof("当前分支%s不支持自动触发构建", h.Ref)

		return fmt.Errorf("the current branch %s does not support automatic triggering of the build", h.Ref)
	}

	buildId, err := jk.Build(job, params)
	if err != nil {
		log.WithFields(log.Fields{
			"event":      h.EventName,
			"project":    h.Project.Name,
			"userName":   h.UserName,
			"params":     params,
			"jenkinsJob": job,
		}).Error(err)
		return err
	}

	log.WithFields(log.Fields{
		"event":    h.EventName,
		"project":  h.Project.Name,
		"userName": h.UserName,
		"params":   params,
	}).Infof("jenkins 构建ID : %d", buildId)

	// TODO: 统计构建结果

	return nil
}
