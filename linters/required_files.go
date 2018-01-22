package linters

import (
	"github.com/simonjohansson/go-linter/model"
	"github.com/spf13/afero"
	"path"
)

type RequiredFilesLinter struct {
	Fs     afero.Fs
	Config model.LinterConfig
}

func (r RequiredFilesLinter) LintManifest(manifest model.Manifest) (model.Result, error) {
	panic("implement me")
}

func (r RequiredFilesLinter) Lint() (model.Result, error) {
	result := model.Result{
		Linter: "Required Files",
	}

	pathToHalfPipeFile := path.Join(r.Config.RepoRoot, ".halfpipe.io")
	exists, err := afero.Exists(r.Fs, pathToHalfPipeFile)
	if err != nil {
		return model.Result{}, err
	}
	if exists == false {
		result.Errors = append(result.Errors, model.Error{
			Message:       "'.halfpipe.io' file is missing",
			Documentation: "",
		})
		return result, nil
	}

	return result, nil
}