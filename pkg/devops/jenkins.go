package devops

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bndr/gojenkins"
	"net/url"
	"path"
	"strconv"
)

type JenkinsOptions struct {
	ctx    context.Context
	Client *gojenkins.Jenkins
}

var Jenkins *JenkinsOptions

func NewJenkins(ctx context.Context, url, username, password string) (*JenkinsOptions, error) {
	jenkins := gojenkins.CreateJenkins(nil, url, username, password)
	_, err := jenkins.Init(ctx)
	if err != nil {
		return nil, err
	}

	return &JenkinsOptions{
		ctx:    ctx,
		Client: jenkins,
	}, nil
}

func (j *JenkinsOptions) Build(name string, params map[string]string) (int64, error) {
	job := gojenkins.Job{Jenkins: j.Client, Raw: new(gojenkins.JobResponse), Base: "/job/" + name}
	endpoint := "/build"
	parameters, err := job.GetParameters(j.ctx)
	if err != nil {
		return 0, err
	}
	if len(parameters) > 0 {
		endpoint = "/buildWithParameters"
	}
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	resp, err := job.Jenkins.Requester.Post(j.ctx, job.Base+endpoint, bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return 0, fmt.Errorf("could not invoke job %q: %s", job.GetName(), resp.Status)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return 0, errors.New("don't have key \"Location\" in response of header")
	}

	u, err := url.Parse(location)
	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseInt(path.Base(u.Path), 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (j *JenkinsOptions) BuildStatus(buildId int64) (string, error) {
	build, err := j.Client.GetBuildFromQueueID(j.ctx, buildId)
	if err != nil {
		return "", err
	}

	return build.GetResult(), nil
}
