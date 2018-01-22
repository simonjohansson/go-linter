package model

type Manifest struct {
	Team  string `yaml:"team"`
	Repo  Repo   `yaml:"repo"`
	Tasks []Task `yaml:"tasks"`
}

type Repo struct {
	Uri        string `yaml:"uri"`
	PrivateKey string `yaml:"private_key"`
}

type Task interface{}

type RunTask struct {
	Script string
	Image  string
	Vars   map[string]string
}

type DockerTask struct {
	Username   string
	Password   string
	Repository string
}

type DeployTask struct {
	Api      string
	Org      string
	Space    string
	Username string
	Password string
	Manifest string
	Vars     map[string]string
}
