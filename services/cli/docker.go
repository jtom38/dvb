package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bitfield/script"
)

const (
	DockerRun                    = "docker run"
	DockerContainerList          = "docker container list"
	DockerContainerInspect       = "docker container inspect"
	DockerContainerStop          = "docker container stop"
	DockerContainerStart         = "docker container start"
	DockerContainerInspectStatus = "docker container inspect -f '{{json .State}}'"

	ContainerStatusStopped = "exited"
	ContainerStatusRunning = "running"

	ErrContainerStopTimeout  = "the requested container did not stop within the requested time frame"
	ErrContainerStartTimeout = "the requested container did not start within the requested time frame"
)

// This client requires the docker cli to be installed.
type DockerCliClient struct{}

func NewDockerCliClient() DockerCliClient {
	return DockerCliClient{}
}

func RunCommand(cmd string) (string, error) {
	var out string

	p := script.NewPipe()
	out, err := p.Exec(cmd).String()
	if err != nil {
		return out, err
	}
	return out, nil
}

func (c DockerCliClient) ListContainers() (string, error) {
	cmd := fmt.Sprintf("%v", DockerContainerList)
	return RunCommand(cmd)
}

func (c DockerCliClient) InspectContainer(name string) (string, error) {
	cmd := fmt.Sprintf("%v %v", DockerContainerInspect, name)
	return RunCommand(cmd)
}

// This sends the stop command but does not wait for it to go offline.
func (c DockerCliClient) StopContainer(name string) (string, error) {
	cmd := fmt.Sprintf("%v %v", DockerContainerStop, name)
	return RunCommand(cmd)
}

func (c DockerCliClient) PollStartContainer(name string) error {
	maxChecks := 30
	counter := 0

	for {
		if counter == maxChecks {
			return errors.New(ErrContainerStopTimeout)
		}

		if c.IsRunning(name) {
			return nil
		}

		res, err := c.StartContainer(name)
		if err != nil {
			return errors.New(res)
		}

		time.Sleep(2 * time.Second)
		counter = counter + 1
	}
}

// This will block the thread till the container has stopped.
// This will wait for a total of 60 seconds, if the container does not stop in time, we error out.
func (c DockerCliClient) PollStopContainer(name string) error {
	maxChecks := 30
	counter := 0

	for {
		if counter == maxChecks {
			return errors.New(ErrContainerStopTimeout)
		}

		if c.IsStopped(name) {
			return nil
		}

		res, err := c.StopContainer(name)
		if err != nil {
			return errors.New(res)
		}

		time.Sleep(2 * time.Second)
		counter = counter + 1
	}
}

func (c DockerCliClient) IsStopped(name string) bool {
	details, _ := c.InspectContainerStatus(name)
	return details.Status == ContainerStatusStopped
}

func (c DockerCliClient) IsRunning(name string) bool {
	details, _ := c.InspectContainerStatus(name)
	return details.Status == ContainerStatusRunning
}

// This sends the start command but does not wait for it to come online.
func (c DockerCliClient) StartContainer(name string) (string, error) {
	cmd := fmt.Sprintf("%v %v", DockerContainerStart, name)
	return RunCommand(cmd)
}

type DockerContainerStatus struct {
	Status     string    `json:"Status"`
	Running    bool      `json:"Running"`
	Paused     bool      `json:"Paused"`
	Restarting bool      `json:"Restarting"`
	OOMKilled  bool      `json:"OOMKilled"`
	Dead       bool      `json:"Dead"`
	Pid        int       `json:"Pid"`
	ExitCode   int       `json:"ExitCode"`
	Error      string    `json:"Error"`
	StartedAt  time.Time `json:"StartedAt"`
	FinishedAt time.Time `json:"FinishedAt"`
}

func (c DockerCliClient) InspectContainerStatus(name string) (DockerContainerStatus, error) {
	var result DockerContainerStatus

	cmd := fmt.Sprintf("%v %v", DockerContainerInspectStatus, name)

	res, err := RunCommand(cmd)
	if err != nil {
		return result, errors.New(res)
	}

	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

type DockerBackupVolumeParams struct {
	ContainerName  string
	BackupFolder   string
	BackupFilename string
	TargetFolder   string
}

func (c DockerCliClient) BackupDockerVolume(params DockerBackupVolumeParams) (string, error) {
	// docker run --rm --volumes-from webdav-app-1 -v $PWD:/backup-dir ubuntu tar cvf /backup-dir/webdav-backup.tar /var/lib/dav

	cmd := fmt.Sprintf("%v --rm --volumes-from %v -v %v:/backup-dir ubuntu tar cvf /backup-dir/%v.tar %v", DockerRun, params.ContainerName, params.BackupFolder, params.BackupFilename, params.TargetFolder)

	return RunCommand(cmd)
}
