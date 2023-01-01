package lib_test

import (
	"fmt"
	"testing"

	"github.com/jtom38/dvb/services/lib"
)

func TestDockerInspectContainerStatus(t *testing.T) {
	client := lib.NewDockerCliClient()
	r, err := client.InspectContainerStatus("webdav-app-1")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("r.Status: %v\n", r.Status)
}

func TestDockerPullStopContainer(t *testing.T) {
	client := lib.NewDockerCliClient()
	err := client.PollStopContainer("webdav-app-1")
	if err != nil {
		t.Error(err)
	}
}

func TestDockerPullStartContainer(t *testing.T) {
	client := lib.NewDockerCliClient()
	err := client.PollStartContainer("webdav-app-1")
	if err != nil {
		t.Error(err)
	}
}
