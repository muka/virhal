package agent

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/project"
)

func TestStart(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	log.Println("Testing deploy")

	project, err := project.NewProjectFromFile("../examples/example1.yml")
	if err != nil {
		t.Fatalf("Failed to load project: %s", err.Error())
		t.FailNow()
	}

	err = Start(project)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}

	status, err = Status(project)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	log.Debugf("Status: %s", status)

	err = Stop(project)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}

}
