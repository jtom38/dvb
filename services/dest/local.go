package dest

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jtom38/dvb/domain"
)

type MoveClient struct {
	backupName    string
	backupPath    string
	containerName string
	destination   string
	fileExtension string
}

func NewMoveClient(backupName, backupPath, containerName, destination string) MoveClient {
	c := MoveClient{
		backupName:    backupName,
		backupPath:    backupPath,
		containerName: containerName,
		destination:   destination,
		fileExtension: "tar",
	}
	return c
}

func NewLocalClient(backupName, backupPath, containerName, destination string) MoveClient {
	c := MoveClient{
		backupName:    backupName,
		backupPath:    backupPath,
		containerName: containerName,
		destination:   destination,
		fileExtension: "tar",
	}
	return c
}

func (c MoveClient) Move(details domain.RunDetails) error {
	// Make sure the file and dest exist
	_, err := os.Stat(c.destination)
	if err != nil {
		return err
	}

	//fileName := fmt.Sprintf("%v.%v", c.backupName, c.fileExtension)

	// Create the directory based on the name of the container/service
	_, err = os.Stat(details.Dest.Local.Directory)
	if err != nil {
		err = os.Mkdir(details.Dest.Local.Directory, 0755)
		if err != nil {
			return err
		}
	}

	err = c.CopyFile(details.Backup.FullFilePath, details.Dest.Local.FullFilePath)
	if err != nil {
		return err
	}

	_, err = os.Stat(details.Dest.Local.FullFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (c MoveClient) CopyFile(source, dest string) error {
	// Check to make sure the source exist
	_, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Open the source file into memory
	s, err := os.Open(source)
	if err != nil {
		return err
	}

	// Check to make sure that the destination does not exist
	// we want an error
	_, err = os.Stat(dest)
	if err == nil {
		return errors.New("destination file already exists, copy job did not start")
	}

	// create the file
	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	s.Close()
	d.Close()

	return nil
}

type RetainClient struct {
	config        domain.ConfigDestLocal
	days          int
	containerName string
	dirPath       string
}

func NewLocalRetainClient(config domain.ConfigDestLocal, containerName string, retainDays int) *RetainClient {
	path := filepath.Join(config.Path, containerName)
	return &RetainClient{
		config:        config,
		days:          retainDays,
		containerName: containerName,
		dirPath: path,
	}
}

func (c RetainClient) GetDirectoryPath() string {
	return c.dirPath
}

func (c RetainClient) Check(pattern string) error {
	if c.config.Path == "" {
		log.Print("Path was empty so skipping Retain check")
		return nil
	}

	if c.days == 0 {
		log.Print("Days was 0 so skipping Retain check")
		return nil
	}

	// Check the number of files
	files, err := c.CountFiles(pattern, c.dirPath)
	if err != nil {
		return err
	}

	if files == 0 {
		log.Print("No files found that matches the pattern.")
		return nil
	}

	if files <= c.days {
		log.Print("Not enough files in the directory to remove ")
		return nil
	}

	// Find the oldest file to remove
	file, err := c.FindOldestFile(pattern, c.dirPath)
	if err != nil {
		return nil
	}

	// build the path to the file based on what we know.
	backupFile := fmt.Sprintf("%v/%v", c.dirPath, file.Name())

	// confirm we have a file the exists
	_, err = os.Stat(backupFile)
	if err != nil {
		return err
	}

	log.Printf("> Removing: %v", backupFile)
	err = os.Remove(backupFile)
	if err != nil {
		return err
	}

	return nil
}

func (c RetainClient) FindOldestFile(pattern, path string) (fs.FileInfo, error) {
	var oldest fs.FileInfo

	files, err := os.ReadDir(path)
	if err != nil {
		return oldest, err
	}

	for _, file := range files {
		name := file.Name()
		if !strings.Contains(name, pattern) {
			continue
		}

		details, err := file.Info()
		if err != nil {
			return oldest, err
		}

		// If we don't have anything stored yet, grab the first one.
		if oldest == nil {
			oldest = details
			continue
		}

		if details.ModTime().Before(oldest.ModTime()) {
			oldest = details
		}
	}

	return oldest, nil
}

func (c RetainClient) CountFiles(pattern, path string) (int, error) {
	found := 0

	dir, err := os.ReadDir(path)
	if err != nil {
		return found, err
	}

	for _, item := range dir {
		if strings.Contains(item.Name(), pattern) {
			found = found + 1
		}
	}

	return found, nil
}
