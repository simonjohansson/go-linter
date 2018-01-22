package linters

import (
	"github.com/simonjohansson/go-linter/model"
	"github.com/spf13/afero"
)

type Linter struct {
	Config model.LinterConfig
	Fs     afero.Fs
}

type LinterI interface {
	Lint() (model.Result, error)
}

func (l Linter) Lint() ([]model.Result, error) {
	results := []model.Result{}
	r, err := NewRequiredFilesLinter(l.Fs, l.Config).Lint()
	if err != nil {
		return []model.Result{}, err
	}
	results = append(results, r)
	if len(r.Errors) != 0 {
		return results, nil
	}

	return results, nil
}
