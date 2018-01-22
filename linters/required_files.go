package linters

import (
	"github.com/simonjohansson/go-linter/model"
	"github.com/spf13/afero"
	"path"
)

type requiredFiles struct {
	fs     afero.Fs
	config model.LinterConfig
}

func NewRequiredFilesLinter(fs afero.Fs, config model.LinterConfig) requiredFiles {
	return requiredFiles{
		fs:     fs,
		config: config,
	}
}

func (r requiredFiles) Lint() (model.Result, error) {
	pathToHalfPipeFile := path.Join(r.config.Path, ".halfpipe.io")
	exists, err := afero.Exists(r.fs, pathToHalfPipeFile)
	if err != nil {
		return model.Result{}, err
	}
	if exists == false {
		return model.Result{
			"Required Files",
			[]model.Error{
				model.Error{
					Message:       "'.halfpipe.io' file is missing",
					Documentation: "",
				},
			},
		}, nil

	}
	return model.Result{}, nil
}
