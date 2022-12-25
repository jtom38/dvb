package domain

import (
	"fmt"
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
	c.Message = append(c.Message, "> ")
	return c
}

func (c *Logs) Add(Message string) {
	log.Print(Message)
	c.Message = append(c.Message, Message)
}

func (c *Logs) Error(err error) {
	log.Print(err.Error())
	c.Message = append(c.Message, err.Error())
}

func (c Logs) generateMessage(msg string) string {
	return fmt.Sprintf("%v", msg)
}
