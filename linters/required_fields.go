package linters

import (
	"github.com/simonjohansson/go-linter/model"
)

type RequiredFieldsLinter struct{}

func (RequiredFieldsLinter) Lint() (model.Result, error) {
	panic("implement me")
}

func (RequiredFieldsLinter) LintManifest(manifest model.Manifest) (model.Result, error) {
	errors := []model.Error{}
	if (manifest.Team == "") {
		errors = append(errors, model.Error{Message: "Required top level field 'team' missing"})
	}
	if (manifest.Repo == model.Repo{}) {
		errors = append(errors, model.Error{Message: "Required top level field 'repo' missing"})
	}
	if (len(manifest.Tasks) == 0) {
		errors = append(errors, model.Error{Message: "Tasks is empty..."})
	}

	return model.Result{
		"Required Fields",
		errors,
	}, nil
}