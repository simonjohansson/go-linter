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

func (r requiredFiles) LintManifest(manifest model.Manifest) (model.Result, error) {
	panic("implement me")
}

func (r requiredFiles) Lint() (model.Result, error) {
	result := model.Result{
		Linter: "Required Files",
	}

	pathToHalfPipeFile := path.Join(r.config.RepoRoot, ".halfpipe.io")
	exists, err := afero.Exists(r.fs, pathToHalfPipeFile)
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

func NewRequiredFilesLinter(fs afero.Fs, config model.LinterConfig) requiredFiles {
	return requiredFiles{
		fs:     fs,
		config: config,
	}
}
