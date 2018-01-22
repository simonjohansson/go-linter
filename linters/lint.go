package linters

import (
	"github.com/simonjohansson/go-linter/model"
	"github.com/simonjohansson/go-linter/manifest"
)

type Linter interface {
	Lint() (model.Result, error)
	LintManifest(manifest model.Manifest) (model.Result, error)
}

type FullLint struct {
	Config               model.LinterConfig
	ManifestReader       manifest.ManifestReader
	RequiredFilesLinter  Linter
	RequiredFieldsLinter Linter
	RepoLinter           Linter
}

func (l FullLint) Lint() ([]model.Result, error) {
	results := []model.Result{}
	r, err := l.RequiredFilesLinter.Lint()
	if err != nil {
		return []model.Result{}, err
	}
	results = append(results, r)
	if len(r.Errors) != 0 {
		// If halfpipe.io file is missing, no need to continue..
		return results, nil
	}

	manifest, err := l.ManifestReader.ParseManifest(l.Config.RepoRoot)
	if err != nil {
		return []model.Result{}, err
	}
	linters := []Linter{
		l.RequiredFieldsLinter,
		l.RepoLinter,
	}
	for _, linter := range linters {
		r, err = linter.LintManifest(manifest)
		if err != nil {
			return []model.Result{}, err
		}
		results = append(results, r)
	}

	return results, nil
}
