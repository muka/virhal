package project

import (
	"context"
	"errors"

	opts "github.com/docker/libcompose/project/options"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/project/options"
)

//Start a project
func (p *Project) Start(opt options.Start) error {

	errs := make([]error, 0)
	for _, service := range p.Services {
		err := service.Start(&opt)
		if err != nil {
			log.Errorf("Failed to start %s: %s", service.Name, err.Error())
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("Failed to start")
	}

	return nil
}

//Start a project
func (s *Service) Start(opt *options.Start) error {

	compose, err := s.GetComposeProject()
	if err != nil {
		return err
	}

	err = compose.Up(context.Background(), opts.Up{})

	if err != nil {
		return err
	}

	return nil
}
