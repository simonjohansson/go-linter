package linters

import (
	"github.com/simonjohansson/go-linter/model"
	"regexp"
)

type RepoLinter struct{}

func (r RepoLinter) isPrivateRepo(repo string) bool {
	regex := `git@github.com:[a-zA-Z0-9]+\/[a-zA-Z0-9_-]+.git`
	matches, _ := regexp.MatchString(regex, repo)
	return matches
}

func (r RepoLinter) isPublicRepo(repo string) bool {
	regex := `https:\/\/github.com\/[a-zA-Z0-9]+\/[a-zA-Z0-9]+`
	matches, _ := regexp.MatchString(regex, repo)
	return matches
}

func (r RepoLinter) LintManifest(manifest model.Manifest) (model.Result, error) {
	result := model.Result{
		Linter: "Repo",
		Errors: []model.Error{},
	}
	if !r.isPrivateRepo(manifest.Repo.Uri) && !r.isPublicRepo(manifest.Repo.Uri) {
		result.Errors = append(result.Errors, model.Error{
			Message: "'" + manifest.Repo.Uri + "' does not look like a real repo!",
		})
	}

	if r.isPrivateRepo(manifest.Repo.Uri) && manifest.Repo.PrivateKey == "" {
		result.Errors = append(result.Errors, model.Error{
			Message: "It looks like you are refering to a private repo, but no private key provided in `repo.private_key`",
		})
	}
	return result, nil
}

func (r RepoLinter) Lint() (model.Result, error) {
	panic("implement me")
}
