package project

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	opts "github.com/muka/virhal/project/options"
)

//Status of a project
func (p *Project) Status(opt opts.Status) (string, error) {

	responses := make([]string, 0)
	errs := make([]error, 0)
	for _, service := range p.Services {
		fmt.Printf("Service %s\n", service.Name)
		res, err := service.Status(&opt)
		if err != nil {
			log.Errorf("Failed to get status of %s: %s", service.Name, err.Error())
			errs = append(errs, err)
		} else {
			responses = append(responses, res)
		}
	}

	if len(errs) > 0 {
		return "", errors.New("Failed to get status")
	}

	return strings.Join(responses, ""), nil
}

//Status of a service
func (s *Service) Status(opt *opts.Status) (string, error) {

	compose, err := s.GetComposeProject()
	if err != nil {
		return "", err
	}

	composeServicesNames := make([]string, 0)
	for _, serviceName := range compose.ServiceConfigs.Keys() {
		composeServicesNames = append(composeServicesNames, serviceName)
	}

	infoset, err := compose.Ps(context.Background(), composeServicesNames...)
	if err != nil {
		return "", err
	}

	if len(infoset) == 0 {
		fmt.Println("Not running")
		return "", nil
	}

	res := "Name\t\t\t\tState\n"
	for _, info := range infoset {
		res += fmt.Sprintf("%s\t\t%s\n", info["Name"], info["State"])
	}

	log.Print(res)
	return res, nil
}
