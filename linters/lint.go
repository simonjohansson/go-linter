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
	Config         model.LinterConfig
	Linters        []Linter
	ManifestReader manifest.ManifestReader
}

func (l FullLint) Lint() ([]model.Result, error) {
	results := []model.Result{}
	// First linter should be required files linter...
	r, err := l.Linters[0].Lint()
	if err != nil {
		return []model.Result{}, err
	}
	results = append(results, r)
	if len(r.Errors) != 0 {
		return results, nil
	}

	manifest, err := l.ManifestReader.ParseManifest(l.Config.RepoRoot)
	for _, linter := range l.Linters[1:] {
		r, err = linter.LintManifest(manifest)
		if err != nil {
			return []model.Result{}, err
		}
		results = append(results, r)
	}

	return results, nil
}
