package docker

import (
	"log"
	"testing"
)

func TestInspectNodes(t *testing.T) {

	nodes, err := SwarmNodesInspect()
	if err != nil {
		log.Fatalf("Error %s", err.Error())
		t.FailNow()
	}

	for _, node := range nodes {
		log.Printf("%v", node.Meta)
		log.Printf("%v", node.ManagerStatus)
		log.Printf("%v", node.Status.Addr)
		log.Printf("%v", node.Spec)
	}

}

func TestSwarmDeployContainer(t *testing.T) {

	err := SwarmDeployContainer()
	if err != nil {
		log.Fatalf("Error %s", err.Error())
		t.FailNow()
	}

}
