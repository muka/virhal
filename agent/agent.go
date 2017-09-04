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

const grpcPort = ":50051"

//Start a project on avail nodes
func Start(project *project.Project) error {
	return eachNode(func(node swarm.Node) error {

		c, conn, err := getClient(node.Status.Addr)
		if err != nil {
			log.Errorf("Cannot get pb client: %s", err.Error())
			return err
		}
		defer conn.Close()

		ctx := context.Background()
		p, err := createPbProject(project)
		if err != nil {
			return err
		}

		response, err := c.Start(ctx, p)
		return checkResponse(response, err)
	})
}

//Stop a project on avail nodes
func Stop(project *project.Project) error {
	return eachNode(func(node swarm.Node) error {

		c, conn, err := getClient(node.Status.Addr)
		if err != nil {
			log.Errorf("Cannot get pb client: %s", err.Error())
			return err
		}
		defer conn.Close()

		ctx := context.Background()
		p, err := createPbProject(project)
		if err != nil {
			log.Errorf("Failed to convert to pb model: %s", err.Error())
			return err
		}

		response, err := c.Stop(ctx, p)
		return checkResponse(response, err)
	})
}

//Status check status for a project on avail nodes
func Status(project *project.Project) (string, error) {
	responses := ""
	err := eachNode(func(node swarm.Node) error {

		responses += fmt.Sprintf("\nNode %s (%s)\n", node.Description.Hostname, node.Status.Addr)
		c, conn, err := getClient(node.Status.Addr)
		if err != nil {
			log.Errorf("Cannot get pb client: %s", err.Error())
			return err
		}
		defer conn.Close()

		ctx := context.Background()
		p, err := createPbProject(project)
		if err != nil {
			log.Errorf("Failed to convert to pb model: %s", err.Error())
			return err
		}

		response, err := c.Status(ctx, p)
		err = checkResponse(response, err)
		if err != nil {
			return err
		}

		responses += response.GetMessage()
		return nil
	})

	return responses, err
}

func eachNode(fn func(node swarm.Node) error) error {

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

		err := fn(node)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Failed execution on nodes: %s", errs[0].Error())
	}

	return nil
}

func checkResponse(response *pb.Response, err error) error {

	if err != nil {
		log.Errorf("Request failed: %s", err.Error())
		return err
	}

	if response.GetCode() >= 400 {
		log.Errorf("Request response has error: %s", err.Error())
		return errors.New("grpc error: " + response.GetMessage())
	}

	return nil
}

func getClient(address string) (pb.ProjectServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address+grpcPort, grpc.WithInsecure())
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
func RunServer() error {

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Errorf("Failed to listen %s: %s", grpcPort, err.Error())
		return err
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
