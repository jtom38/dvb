package domain

import "fmt"

type RunDetails struct {
	*ContainerConfig
	ContainerName       string
	ContainerBackupPath string
	//BackupDirectory     string
	BackupName string
	//BackupExtension     string
	BackupPath string
}

type Logs struct {
	Message []string
}

func (c *Logs) Add(Message string) {
	c.Message = append(c.Message, Message)
}

func (c *Logs) Error(err error) {
	c.Message = append(c.Message, fmt.Sprintf("> Error: %v", err.Error()))
}
