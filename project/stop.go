package project

import (
	"context"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/project/options"
)

//Stop a project
func (p *Project) Stop(opt options.Stop) error {

	errs := make([]error, 0)
	for _, service := range p.Services {
		err := service.Stop(&opt)
		if err != nil {
			log.Errorf("Failed to start %s: %s", service.Name, err.Error())
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("Failed to stop")
	}

	return nil
}

//Stop a project
func (s *Service) Stop(opt *options.Stop) error {

	compose, err := s.GetComposeProject()
	if err != nil {
		return err
	}

	composeServicesNames := make([]string, 0)
	for _, serviceName := range compose.ServiceConfigs.Keys() {
		composeServicesNames = append(composeServicesNames, serviceName)
	}

	err = compose.Kill(context.Background(), "SIGTERM", composeServicesNames...)
	if err != nil {
		return err
	}

	return nil
}
