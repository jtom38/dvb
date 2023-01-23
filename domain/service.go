package domain

import (
	"log"
)

type RunDetails struct {
	ContainerName       string
	ContainerBackupPath string
	Backup              RunBackupDetails
	Dest                RunDestDetails
}

type RunBackupDetails struct {
	TargetDirectory       string
	LocalDirectory        string
	ServiceName           string
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

func NewLogs() Logs {
	c := Logs{}
	return c
}

func (c *Logs) Add(Message string) {
	c.Message = append(c.Message, c.generateMessage(Message))
	log.Print(Message)
}

func (c *Logs) Error(err error) {
	c.Message = append(c.Message, c.generateMessage(err.Error()))
	log.Print(err.Error())
}

func (c Logs) generateMessage(msg string) string {
	return msg
}
