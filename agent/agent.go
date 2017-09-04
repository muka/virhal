package agent

import (
	"context"
	"errors"
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	pb "github.com/muka/virhal/api"
	"github.com/muka/virhal/docker"
	"github.com/muka/virhal/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Deploy a project to available nodes
func Deploy(project *project.Project) error {

	nodes, err := docker.SwarmNodesInspect()
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	for _, node := range nodes {

		if node.Spec.Availability != swarm.NodeAvailabilityActive {
			continue
		}
		if node.Status.State != swarm.NodeStateReady {
			continue
		}
		if node.Status.Addr == "0.0.0.0" {
			continue
		}

		log.Debugf("Node %s %s %s", node.Spec.Availability, node.Status.Addr, node.Description.Hostname)

		c, conn, err := getClient(node.Status.Addr)
		if err != nil {
			errs = append(errs, err)
			log.Errorf("Cannot get pb client: %s", err.Error())
			continue
		}

		ctx := context.Background()
		p, err := createPbProject(project)
		if err != nil {
			errs = append(errs, err)
			log.Errorf("Failed to convert to pb model: %s", err.Error())
			continue
		}

		response, err := c.Start(ctx, p)
		if err != nil {
			errs = append(errs, err)
			log.Errorf("Request failed: %s", err.Error())
		} else {
			if response.GetCode() >= 400 {
				log.Errorf("Request response has error: %s", err.Error())
				errs = append(errs, errors.New("grpc error: "+response.GetMessage()))
			}
		}

		conn.Close()
	}

	if len(errs) > 0 {
		return fmt.Errorf("Failed to deploy: %s", errs[0].Error())
	}

	return nil
}

func getClient(address string) (pb.ProjectServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewProjectServiceClient(conn), conn, nil
}

func createPbProject(project *project.Project) (*pb.Project, error) {

	p := &pb.Project{
		Version:  project.Version,
		Name:     project.Name,
		Services: make(map[string]*pb.Service),
	}

	for name, service := range project.Services {

		tags := make(map[string]string)
		for _, tag := range service.Tags {
			tags[tag.Name] = tag.Value
		}

		_, err := service.GetComposeFileContent()
		if err != nil {
			log.Errorf("Failed to read compose content: %s", err.Error())
			return nil, err
		}

		s := &pb.Service{
			Name: service.Name,
			Mode: service.Mode,
			File: service.FileContent,
			Tags: tags,
		}

		p.Services[name] = s
	}

	return p, nil
}

//RunServer start the gRPC server
func RunServer(port string) error {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProjectServiceServer(s, pb.NewProjectService())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
