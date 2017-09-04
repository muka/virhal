package project

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
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

//GetComposeFileContent read the compose yml file
func (s *Service) GetComposeFileContent() ([]byte, error) {

	if len(s.FileContent) > 0 {
		return s.FileContent, nil
	}

	file := s.File
	if !filepath.IsAbs(file) {
		file = filepath.Join(s.project.context.WorkDir, file)
	}

	s.Context.FullPath = file
	s.Context.WorkDir = filepath.Dir(file)

	log.Debugf("Loading compose file %s", s.Context.FullPath)
	data, err := ioutil.ReadFile(s.Context.FullPath)
	if err != nil {
		return nil, err
	}

	s.FileContent = data

	return s.FileContent, nil
}

//LoadComposeFile read a compose yml file project
func (s *Service) LoadComposeFile() (*prj.Project, error) {

	data, err := s.GetComposeFileContent()
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
