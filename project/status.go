package project

import (
	"context"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	opts "github.com/muka/virhal/project/options"
)

//Status of a project
func (p *Project) Status(opt opts.Status) error {

	errs := make([]error, 0)
	for _, service := range p.Services {
		fmt.Printf("Service %s\n", service.Name)
		err := service.Status(&opt)
		if err != nil {
			log.Errorf("Failed to get status of %s: %s", service.Name, err.Error())
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("Failed to get status")
	}

	return nil
}

//Status of a service
func (s *Service) Status(opt *opts.Status) error {

	compose, err := s.GetComposeProject()
	if err != nil {
		return err
	}

	composeServicesNames := make([]string, 0)
	for _, serviceName := range compose.ServiceConfigs.Keys() {
		composeServicesNames = append(composeServicesNames, serviceName)
	}

	infoset, err := compose.Ps(context.Background(), composeServicesNames...)
	if err != nil {
		return err
	}

	if len(infoset) == 0 {
		fmt.Println("Not running")
		return nil
	}

	fmt.Println("Name\t\t\t\tState")
	for _, info := range infoset {
		fmt.Printf("%s\t\t%s\n", info["Name"], info["State"])
	}

	return nil
}
