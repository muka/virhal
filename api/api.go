package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/project"
	"github.com/muka/virhal/project/options"
)

type projectService struct{}

//NewProjectService create a new project handler server
func NewProjectService() ProjectServiceServer {
	return new(projectService)
}

var defaultResponse = &Response{200, "Ok"}

func (s *projectService) Start(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Starting `%s`", msg.GetName())
	proj, err := convertProject(msg)
	if err != nil {
		return &Response{500, err.Error()}, err
	}
	err = proj.Start(options.Start{})
	if err != nil {
		log.Errorf("Start failed: %s", err.Error())
		return &Response{500, err.Error()}, err
	}
	return defaultResponse, err
}

func (s *projectService) Status(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Status of `%s`", msg.GetName())
	proj, err := convertProject(msg)
	if err != nil {
		return &Response{500, err.Error()}, err
	}
	res, err := proj.Status(options.Status{})
	if err != nil {
		log.Errorf("Status failed: %s", err.Error())
		return &Response{500, err.Error()}, err
	}
	return &Response{200, res}, nil
}

func (s *projectService) Stop(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Stopping `%s`", msg.GetName())
	proj, err := convertProject(msg)
	if err != nil {
		return &Response{500, err.Error()}, err
	}
	err = proj.Stop(options.Stop{})
	if err != nil {
		log.Errorf("Stop failed: %s", err.Error())
		return &Response{500, err.Error()}, err
	}
	return defaultResponse, err
}

func convertProject(msg *Project) (*project.Project, error) {
	p := project.NewProject()

	p.Name = msg.Name
	p.Version = msg.Version

	for name, srvc := range msg.Services {

		s := project.NewService(p)
		s.Name = srvc.Name
		s.FileContent = srvc.GetFile()

		fileName := p.Name + "_" + s.Name
		re := regexp.MustCompile(`([^a-zA-Z0-9_-]*)`)
		fileName = re.ReplaceAllString(fileName, "")

		s.File = filepath.Join(os.TempDir(), fileName)
		log.Debugf("Storing tmp compose file to %s", s.File)

		err := ioutil.WriteFile(s.File, s.FileContent, 0644)
		if err != nil {
			log.Errorf("Failed to store tmp file %s", s.File)
			return nil, err
		}

		p.Services[name] = s
	}

	return p, nil
}
