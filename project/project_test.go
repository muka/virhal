package project

import (
	"log"
	"testing"
)

func TestParseProject(t *testing.T) {

	src := []byte(`
version: "1.2.3"
services:
    example1:
        tags:
            arch: x86
            type: cloud
        file: ./example1.yml
        mode: compose`)

	project := NewProject()
	err := project.Parse(src)
	if err != nil {
		log.Fatal("Failed to parse: %", err.Error())
		t.FailNow()
	}

	if project.Version != "1.2.3" {
		t.FailNow()
	}

	if _, ok := project.Services["example1"]; !ok {
		t.FailNow()
	}

	if len(project.Services["example1"].Tags) != 2 {
		t.FailNow()
	}

}
