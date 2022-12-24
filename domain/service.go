package domain

import "fmt"

type RunDetails struct {
	ContainerName       string
	ContainerBackupPath string
	Backup              RunBackupDetails
	Dest                RunDestDetails
}

type RunBackupDetails struct {
	Directory             string
	FileName              string
	Extension             string
	FileNameWithExtension string
	FullFilePath          string
}

type RunDestDetails struct {
	Local RunDetailsDestLocal
}

// This struct contains the information on where the local data
type RunDetailsDestLocal struct {
	Directory             string
	FileName              string
	Extension             string
	FileNameWithExtension string
	FullFilePath          string
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
