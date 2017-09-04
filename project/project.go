package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	compose "github.com/docker/libcompose/project"

	"gopkg.in/yaml.v2"
)

//VERSION default project version
const VERSION = "1"

const (
	//DeplyomentModeCompose compose deployment mode
	DeplyomentModeCompose = "compose"
	//DeplyomentModeSwarm swarm deployment mode
	DeplyomentModeSwarm = "swarm"
	//DeplyomentModeFunction function deployment mode
	DeplyomentModeFunction = "function"
)

//YamlService yaml model
type YamlService struct {
	Tags map[string]string `yaml:"tags"`
	File string            `yaml:"file"`
	Mode string            `yaml:"mode"`
}

//YamlProject yaml model
type YamlProject struct {
	Name     string                 `yaml:"name"`
	Version  string                 `yaml:"version"`
	Services map[string]YamlService `yaml:"services"`
}

//Tag definition
type Tag struct {
	Name  string
	Value string
}

//Service definition
type Service struct {
	Context        *Context
	project        *Project
	composeProject *compose.Project

	Name        string
	Tags        map[string]Tag
	File        string
	FileContent []byte
	Mode        string
}

//Context definition
type Context struct {
	WorkDir  string
	FullPath string
}

//Project definition
type Project struct {
	context  *Context
	Name     string
	Version  string
	Services map[string]*Service
}

//NewProject initialize a new project
func NewProject() *Project {
	p := Project{
		context:  &Context{},
		Version:  VERSION,
		Services: make(map[string]*Service),
	}
	return &p
}

//NewService initialize a new service
func NewService(project *Project) *Service {
	p := Service{
		project: project,
		Context: &Context{},
		Mode:    "compose",
		Tags:    make(map[string]Tag),
	}
	return &p
}

//NewProjectFromFile initialize a new project from a configuration file
func NewProjectFromFile(file string) (*Project, error) {

	if file == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		file = dir + "/virhal.yml"
	}

	log.Debugf("Loading %s", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := NewProject()

	configFilePath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	p.context.FullPath = configFilePath
	p.context.WorkDir = filepath.Dir(configFilePath)

	if p.Name == "" {
		bname := filepath.Base(configFilePath)
		bname = strings.Replace(bname, ".yml", "", -1)
		re := regexp.MustCompile(`([^a-zA-Z0-9_-]*)`)
		bname = re.ReplaceAllString(bname, "")
		p.Name = bname
	}

	err = p.Parse(data)
	if err != nil {
		return nil, err
	}
	return p, nil
}

//Parse parse a project content
func (p *Project) Parse(source []byte) error {

	yml := YamlProject{}

	err := yaml.Unmarshal(source, &yml)
	if err != nil {
		return err
	}

	mergeYamlProject(&yml, p)

	return nil
}

func mergeYamlProject(raw *YamlProject, p *Project) {

	if raw.Version != "" {
		p.Version = raw.Version
	}

	for serviceName, rawService := range raw.Services {

		s := Service{
			Context:     new(Context),
			project:     p,
			Name:        serviceName,
			File:        rawService.File,
			Mode:        rawService.Mode,
			Tags:        make(map[string]Tag),
			FileContent: make([]byte, 0),
		}

		for tagname, tagval := range rawService.Tags {
			s.Tags[tagname] = Tag{tagname, tagval}
		}

		p.Services[serviceName] = &s
	}

}
