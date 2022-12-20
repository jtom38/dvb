package domain

import "fmt"

type RunDetails struct {
	//*ContainerConfig
	ContainerName       string
	ContainerBackupPath string

	BackupDirectory string
	BackupFileName  string
	BackupExtension string
	BackupPath      string
}

type Logs struct {
	Message []string
}

func (c *Logs) Add(Message string) {
	c.Message = append(c.Message, c.generateMessage(Message))
}

func (c *Logs) Error(err error) {
	c.Message = append(c.Message, c.generateMessage(err.Error()))
}

func (c Logs) generateMessage(msg string) string {
	return fmt.Sprintf("> %v", msg)
}
