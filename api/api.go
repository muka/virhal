package api

import (
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
)

type projectService struct{}

//NewProjectService create a new project handler server
func NewProjectService() ProjectServiceServer {
	return new(projectService)
}

var defaultResponse = &Response{200, "Ok"}

func (s *projectService) Start(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Start %s (v%s)", msg.GetName(), msg.GetVersion())
	return defaultResponse, nil
}

func (s *projectService) Status(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Status of %s (v%s)", msg.GetName(), msg.GetVersion())
	return defaultResponse, nil
}

func (s *projectService) Stop(ctx context.Context, msg *Project) (*Response, error) {
	log.Debugf("Stop %s (v%s)", msg.GetName(), msg.GetVersion())
	return defaultResponse, nil
}
