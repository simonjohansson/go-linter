package manifest

import (
	"github.com/spf13/afero"
	"github.com/simonjohansson/go-linter/model"
	"github.com/pkg/errors"
	"path"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type ManifestReader interface {
	ParseManifest(repoRoot string) (model.Manifest, error)
}

type manifestReader struct {
	fs afero.Fs
}

func NewManifestReader(fs afero.Fs) ManifestReader {
	return manifestReader{
		fs: fs,
	}
}

func (m manifestReader) parseVars(vars interface{}) map[string]string {
	returnMap := map[string]string{}
	switch vars.(type) {
	case map[interface{}]interface{}:

		for k, v := range vars.(map[interface{}]interface{}) {
			switch v.(type) {
			case string:
				returnMap[k.(string)] = v.(string)
			case int:
				returnMap[k.(string)] = strconv.Itoa(v.(int))
			}

		}
	}
	return returnMap
}

func (m manifestReader) parseRunTask(task map[interface{}]interface{}) model.RunTask {
	runTask := model.RunTask{}
	mapstructure.Decode(task, &runTask)
	if runTask.Vars != nil {
		if vars, ok := task["vars"]; ok {
			runTask.Vars = m.parseVars(vars)
		}
	} else {
		runTask.Vars = make(map[string]string)
	}
	return runTask
}

func (m manifestReader) parseDeployTask(task map[interface{}]interface{}, team string) model.DeployTask {
	deployTask := model.DeployTask{}
	mapstructure.Decode(task, &deployTask)
	if deployTask.Org == "" {
		deployTask.Org = team
	}
	if deployTask.Password == "" {
		deployTask.Password = "((cf-credentials.password))"
	}
	if deployTask.Username == "" {
		deployTask.Username = "((cf-credentials.username))"
	}
	if deployTask.Manifest == "" {
		deployTask.Manifest = "manifest.yml"
	}
	if deployTask.Vars != nil {
		if vars, ok := task["vars"]; ok {
			deployTask.Vars = m.parseVars(vars)
		}
	} else {
		deployTask.Vars = make(map[string]string)
	}
	return deployTask
}

func (m manifestReader) parseDockerTask(task map[interface{}]interface{}) model.DockerTask {
	dockerTask := model.DockerTask{}
	mapstructure.Decode(task, &dockerTask)
	return dockerTask
}

func (m manifestReader) ParseManifest(repoRoot string) (model.Manifest, error) {
	halfPipePath := path.Join(repoRoot, ".halfpipe.io")
	manifest := model.Manifest{}
	exists, err := afero.Exists(m.fs, halfPipePath)
	if err != nil {
		return manifest, err
	}
	if exists == false {
		return manifest, errors.New("Manifest at path '" + halfPipePath + "' does not exist")
	}

	bytes, err := afero.ReadFile(m.fs, halfPipePath)
	if err != nil {
		return manifest, err
	}

	yaml.Unmarshal(bytes, &manifest)
	tasks := []model.Task{}
	for _, task := range manifest.Tasks {
		t := task.(map[interface{}]interface{})
		switch(t["task"]) {
		case "run":
			tasks = append(tasks, m.parseRunTask(t))
		case "docker":
			tasks = append(tasks, m.parseDockerTask(t))
		case "deploy":
			tasks = append(tasks, m.parseDeployTask(t, manifest.Team))
		}
	}

	if len(tasks) > 0 {
		manifest.Tasks = tasks
	}
	return manifest, nil
}
