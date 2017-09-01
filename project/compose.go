package project

import (
	"io/ioutil"
	"path/filepath"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	prj "github.com/docker/libcompose/project"
)

//GetComposeProject return a project instance, loading if needed
func (s *Service) GetComposeProject() (*prj.Project, error) {
	if s.composeProject == nil {
		return s.LoadComposeFile()
	}
	return s.composeProject, nil
}

//LoadComposeFile read a compose yml file project
func (s *Service) LoadComposeFile() (*prj.Project, error) {

	file := s.File

	if !filepath.IsAbs(file) {
		file = filepath.Join(s.project.context.WorkDir, file)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	compose, err := docker.NewProject(&ctx.Context{
		Context: prj.Context{
			ComposeBytes: [][]byte{data},
			ProjectName:  s.project.Name,
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	s.composeProject = compose.(*prj.Project)
	return s.composeProject, nil
}
