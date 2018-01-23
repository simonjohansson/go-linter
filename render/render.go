package render

import (
	"github.com/simonjohansson/go-linter/model"
	"github.com/concourse/atc"
	"fmt"
	"strings"
	"gopkg.in/yaml.v2"
)

type Render interface {
	Render(model.Manifest)
}

type ConcourseRenderer struct {
}

func (ConcourseRenderer) makeGitConfig(repo model.Repo) atc.ResourceConfig {
	source := atc.Source{
		"uri": repo.Uri,
	}
	if repo.PrivateKey != "" {
		source["private_key"] = repo.PrivateKey
	}
	return atc.ResourceConfig{
		Name:   repo.RepoName(),
		Type:   "git",
		Source: source,
	}
}

func (ConcourseRenderer) dockerImageAndTag(image string) (string, string) {
	if strings.Contains(image, ":") {
		split := strings.Split(image, ":")
		return split[0], split[1]
	}
	return image, "latest"
}

func (c ConcourseRenderer) makeRunJob(task model.RunTask, repo model.Repo) atc.JobConfig {
	image, tag := c.dockerImageAndTag(task.Image)
	return atc.JobConfig{
		Name:   task.Script,
		Serial: true,
		Plan: atc.PlanSequence{
			atc.PlanConfig{Get: repo.RepoName(), Trigger: true},
			atc.PlanConfig{
				Task: task.Script,
				TaskConfig: &atc.TaskConfig{
					Platform: "linux",
					Params:   task.Vars,
					ImageResource: &atc.ImageResource{
						Type: "docker-image",
						Source: atc.Source{
							"repository": image,
							"tag":        tag,
						},
					},
					Run: atc.TaskRunConfig{
						Path: "/bin/sh",
						Dir:  repo.RepoName(),
						Args: []string{"-exc", fmt.Sprintf("./%s", task.Script)},
					},
					Inputs: []atc.TaskInputConfig{
						atc.TaskInputConfig{Name: repo.RepoName()},
					},
				}}}}
}

func (c ConcourseRenderer) Render(manifest model.Manifest) atc.Config {
	config := atc.Config{}
	config.Resources = append(config.Resources, c.makeGitConfig(manifest.Repo))
	for _, task := range manifest.Tasks {
		switch task.(type) {
		case model.RunTask:
			config.Jobs = append(config.Jobs, c.makeRunJob(task.(model.RunTask), manifest.Repo))
		}
	}
	return config
}

func (c ConcourseRenderer) RenderToString(manifest model.Manifest) string {
	pipeline := c.Render(manifest)
	renderedPipeline, _ := yaml.Marshal(pipeline)
	return string(renderedPipeline)
}
