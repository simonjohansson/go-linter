package render_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
. "gopkg.in/check.v1"

	"testing"
	. "github.com/simonjohansson/go-linter/render"
	"github.com/simonjohansson/go-linter/model"
	"github.com/concourse/atc"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
)

var _ = Describe("ConcourseRenderer", func() {

	Context("Repo", func() {
		It("renders http repo to a git resource", func() {

			manifest := model.Manifest{
				Repo: model.Repo{
					Uri: "https://github.com/cloudfoundry/bosh-cli",
				},
			}
			pipeline := ConcourseRenderer{}.RenderToString(manifest)
			expected := `
			groups: []
			resources:
			- name: bosh-cli
			  type: git
			  source:
				uri: https://github.com/cloudfoundry/bosh-cli
			resource_types: []
			jobs: []`

			Expect(pipeline).To(Equal(expected))
		})

		It("renders ssh repo to a git resource", func() {
			manifest := model.Manifest{
				Repo: model.Repo{
					Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
					PrivateKey: "((something.secret))",
				},
			}
			pipeline := ConcourseRenderer{}.RenderToString(manifest)
			expected := `
				groups: []
				resources:
				- name: ee-half-pipe-landing
				  type: git
				  source:
					private_key: "((someting.secret))"
					uri: git@github.com:springernature/ee-half-pipe-landing.git
				resource_types: []
				jobs: []`
			Expect(pipeline).To(Equal(expected))
		})
	})

	Context("Tasks", func() {
		Context("run", func() {
			It("Renders a task without vars correctly", func() {
				manifest := model.Manifest{
					Repo: model.Repo{
						Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
						PrivateKey: "((something.secret))",
					},
					Tasks: []model.Task{
						model.RunTask{
							Script: "yolo.sh",
							Image:  "something",
						},
					},
				}

				pipeline := ConcourseRenderer{}.Render(manifest)
				Expect(pipeline.Jobs[0]).To(Equal(atc.JobConfig{
					Name:   "yolo.sh",
					Serial: true,
					Plan: atc.PlanSequence{
						atc.PlanConfig{Get: manifest.Repo.RepoName(), Trigger: true},
						atc.PlanConfig{Task: "yolo.sh", TaskConfig: &atc.TaskConfig{
							Platform: "linux",
							ImageResource: &atc.ImageResource{
								Type: "docker-image",
								Source: atc.Source{
									"repository": "something",
									"tag":        "latest",
								},
							},
							Params: nil,
							Run: atc.TaskRunConfig{
								Path: "/bin/sh",
								Dir:  manifest.Repo.RepoName(),
								Args: []string{"-exc", "./yolo.sh"},
							},
							Inputs: []atc.TaskInputConfig{
								atc.TaskInputConfig{Name: manifest.Repo.RepoName()},
							},
						}},
					}}))
			})

			It("Renders a task with vars correctly", func() {
				manifest := model.Manifest{
					Repo: model.Repo{
						Uri:        "git@github.com:springernature/ee-half-pipe-landing.git",
						PrivateKey: "((something.secret))",
					},
					Tasks: []model.Task{
						model.RunTask{
							Script: "yolo.sh",
							Image:  "something",
							Vars: map[string]string{
								"VAR1": "Value",
								"VAR2": "Value",
							},
						},
					},
				}

				pipeline := ConcourseRenderer{}.Render(manifest)
				Expect(pipeline.Jobs[0]).To(Equal(atc.JobConfig{
					Name:   "yolo.sh",
					Serial: true,
					Plan: atc.PlanSequence{
						atc.PlanConfig{Get: manifest.Repo.RepoName(), Trigger: true},
						atc.PlanConfig{Task: "yolo.sh", TaskConfig: &atc.TaskConfig{
							Platform: "linux",
							Params: map[string]string{
								"VAR1": "Value",
								"VAR2": "Value",
							},
							ImageResource: &atc.ImageResource{
								Type: "docker-image",
								Source: atc.Source{
									"repository": "something",
									"tag":        "latest",
								},
							},
							Run: atc.TaskRunConfig{
								Path: "/bin/sh",
								Dir:  manifest.Repo.RepoName(),
								Args: []string{"-exc", fmt.Sprint("./yolo.sh")},
							},
							Inputs: []atc.TaskInputConfig{
								atc.TaskInputConfig{Name: manifest.Repo.RepoName()},
							},
						}},
					}}))
			})

		})
	})
})
