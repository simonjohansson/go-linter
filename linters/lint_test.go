package linters_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/simonjohansson/go-linter/model"
	. "github.com/simonjohansson/go-linter/linters"
	"github.com/simonjohansson/go-linter/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lint", func() {
	var (
		config   model.LinterConfig
		fullLint FullLint

		// Mocks
		mockCtrl            *gomock.Controller
		requiredFilesLinter *mocks.MockLinter
		requredFieldsLinter *mocks.MockLinter
		repoLinter          *mocks.MockLinter
		manifestReader      *mocks.MockManifestReader
	)

	BeforeEach(func() {
		config = model.LinterConfig{
			RepoRoot: "/path/to/repo",
		}

		mockCtrl = gomock.NewController(GinkgoT())
		requiredFilesLinter = mocks.NewMockLinter(mockCtrl)
		requredFieldsLinter = mocks.NewMockLinter(mockCtrl)
		repoLinter = mocks.NewMockLinter(mockCtrl)
		manifestReader = mocks.NewMockManifestReader(mockCtrl)

		fullLint = FullLint{Config: config,
			ManifestReader: manifestReader,
			RequiredFilesLinter: requiredFilesLinter,
			RequiredFieldsLinter: requredFieldsLinter,
			RepoLinter: repoLinter,
		}
	})

	It("Exits if required files returns errors", func() {
		result := model.Result{Linter: "Blah", Errors: []model.Error{model.Error{}}}
		requiredFilesLinter.EXPECT().Lint().Return(result, nil)

		results, err := fullLint.Lint()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(results)).To(Equal(1))
		Expect(results[0]).To(Equal(result))
	})

	It("Runs all linters", func() {
		result1 := model.Result{Linter: "a"}
		result2 := model.Result{Linter: "b"}
		result3 := model.Result{Linter: "c"}

		requiredFilesLinter.EXPECT().Lint().Return(result1, nil)

		manifest := model.Manifest{}
		manifestReader.EXPECT().ParseManifest(config.RepoRoot).Return(manifest, nil)
		requredFieldsLinter.EXPECT().LintManifest(manifest).Return(result2, nil)
		repoLinter.EXPECT().LintManifest(manifest).Return(result3, nil)

		results, err := fullLint.Lint()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(results)).To(Equal(3))
		Expect(results[0]).To(Equal(result1))
		Expect(results[1]).To(Equal(result2))
		Expect(results[2]).To(Equal(result3))
	})
})
