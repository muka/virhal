package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	compose "github.com/docker/libcompose/project"

	"gopkg.in/yaml.v2"
)

const composeFileContent = "version: \"2\"\nservices:\n  %s:\n    %s"

//VERSION default project version
const VERSION = "1"

const (
	//DeplyomentModeCompose compose deployment mode
	DeplyomentModeCompose = "compose"
	//DeplyomentModeSwarm swarm deployment mode
	DeplyomentModeSwarm = "swarm"
	//DeplyomentModeFunction function deployment mode
	DeplyomentModeFunction = "function"
	//DeplyomentModeContainer single container deployment mode
	DeplyomentModeContainer = "container"
)

//YamlService yaml model
type YamlService struct {
	Tags    map[string]string      `yaml:"tags"`
	File    string                 `yaml:"file"`
	Mode    string                 `yaml:"mode"`
	Service map[string]interface{} `yaml:"service"`
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
		Mode:    DeplyomentModeCompose,
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
	log.Debugf("Project name: %s", p.Name)

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

	err = mergeYamlProject(&yml, p)
	if err != nil {
		return err
	}

	return nil
}

func mergeYamlProject(raw *YamlProject, p *Project) error {

	if raw.Version != "" {
		p.Version = raw.Version
	}

	for serviceName, rawService := range raw.Services {

		deployMode := rawService.Mode
		filePath := rawService.File
		fileContent := make([]byte, 0)

		if rawService.Mode == DeplyomentModeContainer {

			if rawService.File == "" && len(rawService.Service) == 0 {
				return fmt.Errorf("Service %s should have a `file` or `service` field", serviceName)
			}

			deployMode = DeplyomentModeCompose

			//file reference, load content
			if rawService.File != "" {
				content, err := ioutil.ReadFile(rawService.File)
				if err != nil {
					log.Errorf("Failed to read file %s", rawService.File)
					return err
				}

				err = yaml.Unmarshal(content, &rawService.Service)
				if err != nil {
					return err
				}

			}

			//service description, transform to docker-compose file format with single service
			if len(rawService.Service) != 0 {

				serviceContent, err := yaml.Marshal(rawService.Service)
				if err != nil {
					return err
				}

				fileName := strings.Join([]string{p.Name, serviceName}, "-") + ".yml"
				filePath = filepath.Join(os.TempDir(), fileName)
				fileContent = []byte(fmt.Sprintf(
					composeFileContent,
					serviceName,
					serviceContent,
				))

				log.Debugf("Storing tmp compose file to %s", filePath)
				err = ioutil.WriteFile(filePath, fileContent, 0644)
				if err != nil {
					log.Errorf("Failed to store tmp file %s", filePath)
					return err
				}
			}
		}

		s := Service{
			Context:     new(Context),
			project:     p,
			Name:        serviceName,
			File:        filePath,
			Mode:        deployMode,
			Tags:        make(map[string]Tag),
			FileContent: fileContent,
		}

		for tagname, tagval := range rawService.Tags {
			s.Tags[tagname] = Tag{tagname, tagval}
		}

		p.Services[serviceName] = &s
	}

	return nil
}
