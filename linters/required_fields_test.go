package linters_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/simonjohansson/go-linter/model"
	"github.com/simonjohansson/go-linter/linters"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequiredFields", func() {

	It("Returns error if team is missing", func() {
		manifest := model.Manifest{
			Repo:  model.Repo{Uri: "asdf"},
			Tasks: []model.Task{model.RunTask{}},
		}
		results, _ := linters.RequiredFieldsLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if repo is missing", func() {
		manifest := model.Manifest{
			Team:  "asdf",
			Tasks: []model.Task{model.RunTask{}},
		}
		results, _ := linters.RequiredFieldsLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if tasks is missing", func() {
		manifest := model.Manifest{
			Team: "asdf",
			Repo: model.Repo{Uri: "asdf"},
		}
		results, _ := linters.RequiredFieldsLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})
})
