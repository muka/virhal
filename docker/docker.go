package docker

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/muka/virhal/emitter"
)

var dockerClient *client.Client

//Event wrapper event for emitter
type Event struct {
	Message events.Message
}

//GetRaw return message as interface
func (ev *Event) GetRaw() interface{} {
	return ev.Message
}

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
func SwarmDeployContainer() error {
	_, err := GetEnvClient()
	if err != nil {
		return err
	}

	return nil
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
			emitter.GetEmitter().Emit(msg.Type, Event{msg})
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
