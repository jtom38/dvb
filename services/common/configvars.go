package common

import (
	"os"
	"strings"
	"time"
)

const (
	ConfigVariablePwd     = "{{PWD}}"
	ConfigVariableDate    = "{{DATE}}"
	ConfigVariableUserDir = "{{USERDIR}}"
)

// This function finds all the known config variables and replaces the values.
func ReplaceAllConfigVariables(value string) (string, error) {
	var err error
	var dir string
	t := value

	if strings.Contains(value, ConfigVariablePwd) {
		dir, err = os.Getwd()
		if err != nil {
			return t, err
		}
		t = strings.ReplaceAll(t, ConfigVariablePwd, dir)
	}

	if strings.Contains(value, ConfigVariableDate) {
		todayString := time.Now().Format("20060102")
		t = strings.ReplaceAll(t, ConfigVariableDate, todayString)
	}

	if strings.Contains(value, ConfigVariableUserDir) {
		dir, err = os.UserHomeDir()
		t = strings.ReplaceAll(t, ConfigVariableUserDir, dir)
	}

	return t, nil
}
