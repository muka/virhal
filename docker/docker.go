package docker

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

//GetEnvClient return the environment client instance
func GetEnvClient() (*client.Client, error) {

	var err error
	if dockerClient == nil {
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			return nil, err
		}
	}

	return dockerClient, nil
}

//SwarmNodesInspect return the available nodes metadata
func SwarmNodesInspect() ([]swarm.Node, error) {
	c, err := GetEnvClient()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	listOpts := types.NodeListOptions{}
	nodes, err := c.NodeList(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	return nodes, err
}

//SwarmDeployContainer return the available nodes metadata
func SwarmDeployContainer() ([]swarm.Node, error) {
	c, err := GetEnvClient()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	listOpts := types.NodeListOptions{}
	nodes, err := c.NodeList(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	return nodes, err
}

// WatchEvents watch for docker events
func WatchEvents() error {

	c, err := GetEnvClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	opts := types.EventsOptions{}
	msgChan, errChan := c.Events(ctx, opts)
	quit := false
	for {
		select {
		case msg := <-msgChan:
			log.Debugf("Emitting event %s", msg.Type)
			GetEmitter().Emit(msg.Type, Event{msg})
		case err := <-errChan:
			if err != nil {
				fmt.Printf("got error %v\n", err.Error())
				quit = true
				break
			}
		}
		if quit {
			break
		}
	}

	cancel()
	// call again to reopen the API stream after an error
	return WatchEvents()
}
