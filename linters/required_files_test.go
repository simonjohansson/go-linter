package linters_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/spf13/afero"
	"github.com/simonjohansson/go-linter/model"
	"github.com/simonjohansson/go-linter/linters"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequiredFiles", func() {
	var (
		fs     afero.Fs
		config model.LinterConfig
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		config = model.LinterConfig{
			RepoRoot: "/path/to/repo",
		}
	})

	It("Returns error if .halfpipe.io is missing", func() {
		results, _ := linters.RequiredFilesLinter{fs, config}.Lint()
		Expect(len(results.Errors)).To(Equal(1))
	})

	It("Returns empty error if .halfpipe.io is present", func() {
		afero.WriteFile(fs, "/path/to/repo/.halfpipe.io", []byte(""), 0644)
		results, _ := linters.RequiredFilesLinter{fs, config}.Lint()
		Expect(len(results.Errors)).To(Equal(0))
	})
})
