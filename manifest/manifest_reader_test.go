package manifest_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/spf13/afero"
	. "github.com/onsi/gomega"
	"github.com/simonjohansson/go-linter/manifest"
	"github.com/simonjohansson/go-linter/model"
)

var _ = Describe("RequiredFiles", func() {
	var (
		fs afero.Fs
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
	})

	It("Returns error if .halfpipe.io is missing", func() {
		manifestReader := manifest.NewManifestReader(fs)

		_, err := manifestReader.ParseManifest("/path/to/repo")
		Expect(err).To(HaveOccurred())
	})

	It("Returns empty manifest is .halfpipe.io is empty", func() {
		manifestReader := manifest.NewManifestReader(fs)
		afero.WriteFile(fs, "/path/to/repo/.halfpipe.io", []byte(""), 0644)

		manifest, err := manifestReader.ParseManifest("/path/to/repo")
		Expect(err).To(Not(HaveOccurred()))
		Expect(manifest).To(Equal(model.Manifest{}))
	})

	It("Parses empty .halfpipe.io to empty manifest", func() {
		manifestReader := manifest.NewManifestReader(fs)
		content := ``
		afero.WriteFile(fs, "/path/to/repo/.halfpipe.io", []byte(content), 0644)

		manifest, err := manifestReader.ParseManifest("/path/to/repo")
		Expect(err).To(Not(HaveOccurred()))
		Expect(manifest).To(Equal(model.Manifest{}))
	})

	It("Parses minimal .halfpipe.io to minimal manifest", func() {
		manifestReader := manifest.NewManifestReader(fs)
		content := `
team: engineering-enablement
`
		afero.WriteFile(fs, "/path/to/repo/.halfpipe.io", []byte(content), 0644)

		manifest, err := manifestReader.ParseManifest("/path/to/repo")
		Expect(err).To(Not(HaveOccurred()))
		Expect(manifest).To(Equal(model.Manifest{
			Team: "engineering-enablement",
		}))
	})

	It("Parses .halfpipe.io to manifest", func() {
		manifestReader := manifest.NewManifestReader(fs)
		content := `
team: engineering-enablement
repo:
  uri: https://....
  private_key: asdf
tasks:
- task: run
  script: ./test.sh
  image: openjdk:8-slim
- task: docker
  username: ((docker.username))
  password: ((docker.password))
  repository: simonjohansson/half-pipe-linter 
- task: deploy
  space: test
  api: https://api.europe-west1.cf.gcp.springernature.io
- task: run
  script: ./asdf.sh
  image: openjdk:8-slim
  vars:
    A: asdf
    B: 1234
- task: deploy
  space: test
  api: https://api.europe-west1.cf.gcp.springernature.io
  vars:
    VAR1: asdf1234
    VAR2: 9876
`
		afero.WriteFile(fs, "/path/to/repo/.halfpipe.io", []byte(content), 0644)

		manifest, err := manifestReader.ParseManifest("/path/to/repo")
		Expect(err).To(Not(HaveOccurred()))
		Expect(manifest).To(Equal(model.Manifest{
			Team: "engineering-enablement",
			Repo: model.Repo{
				Uri:        "https://....",
				PrivateKey: "asdf",
			},
			Tasks: []model.Task{
				model.RunTask{
					Script: "./test.sh",
					Image:  "openjdk:8-slim",
					Vars:   make(map[string]string),
				},
				model.DockerTask{
					Username:   "((docker.username))",
					Password:   "((docker.password))",
					Repository: "simonjohansson/half-pipe-linter",
				},
				model.DeployTask{
					Username: "((cf-credentials.username))",
					Password: "((cf-credentials.password))",
					Api:      "https://api.europe-west1.cf.gcp.springernature.io",
					Org:      "engineering-enablement",
					Space:    "test",
					Manifest: "manifest.yml",
					Vars:     make(map[string]string),
				},
				model.RunTask{
					Script: "./asdf.sh",
					Image:  "openjdk:8-slim",
					Vars: map[string]string{
						"A": "asdf",
						"B": "1234",
					},
				},
				model.DeployTask{
					Username: "((cf-credentials.username))",
					Password: "((cf-credentials.password))",
					Api:      "https://api.europe-west1.cf.gcp.springernature.io",
					Org:      "engineering-enablement",
					Space:    "test",
					Manifest: "manifest.yml",
					Vars: map[string]string{
						"VAR1": "asdf1234",
						"VAR2": "9876",
					},
				},
			},
		}))
	})

})
