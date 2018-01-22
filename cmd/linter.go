package main

import (
	"flag"
	"github.com/simonjohansson/go-linter/model"
	"os"
	"path/filepath"
	"github.com/simonjohansson/go-linter/linters"
	"fmt"
	"syscall"
	"github.com/spf13/afero"
	"github.com/simonjohansson/go-linter/manifest"
)

func getPath(path string) (string, error) {
	if path == "." || path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", nil
	}
	return path, nil
}

func config() (model.LinterConfig, error) {
	vaultToken := flag.String("vaultToken", "", "Vault token")
	path := flag.String("path", "", "ath to repo to lint, if not specified pwd will be linted")

	flag.Parse()
	p, err := getPath(*path)
	if err != nil {
		return model.LinterConfig{}, err
	}

	return model.LinterConfig{
		VaultToken: *vaultToken,
		RepoRoot:   p,
	}, nil
}

func main() {
	config, err := config()
	if err != nil {
		fmt.Println(err)
		syscall.Exit(-1)
	}

	results, err := linters.FullLint{
		Config:               config,
		ManifestReader:       manifest.NewManifestReader(afero.NewOsFs()),
		RequiredFilesLinter:  linters.RequiredFilesLinter{afero.NewOsFs(), config},
		RequiredFieldsLinter: linters.RequiredFieldsLinter{},
		RepoLinter:           linters.RepoLinter{},
	}.Lint()

	if err != nil {
		fmt.Println(err)
		syscall.Exit(-1)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
