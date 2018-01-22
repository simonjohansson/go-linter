package linters_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/simonjohansson/go-linter/model"
	"github.com/simonjohansson/go-linter/linters"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {

	It("Returns error if repo.uri doesnt look like git repo", func() {
		manifest := model.Manifest{
			Repo: model.Repo{
				Uri: "gitasdasdt@github.com:simonjohansson/go-linter.git",
			},
		}
		results, _ := linters.RepoLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))

		manifest = model.Manifest{
			Repo: model.Repo{
				Uri: "https://github.coasdm/cenkalti/backoff",
			},
		}
		results, _ = linters.RepoLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns error if repo.uri is private but no private key specified", func() {
		manifest := model.Manifest{
			Repo: model.Repo{
				Uri: "git@github.com:simonjohansson/go-linter.git",
			},
		}
		results, _ := linters.RepoLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns no errors if all is in order", func() {
		manifest := model.Manifest{
			Repo: model.Repo{
				Uri: "git@github.com:simonjohansson/go-linter.git",
				PrivateKey: "asd",
			},
		}
		results, _ := linters.RepoLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(0))

		manifest = model.Manifest{
			Repo: model.Repo{
				Uri: "https://github.com/cenkalti/backoff",
			},
		}
		results, _ = linters.RepoLinter{}.LintManifest(manifest)
		Expect(len(results.Errors)).To(Equal(0))
	})
})
