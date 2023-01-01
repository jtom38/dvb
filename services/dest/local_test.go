package dest_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/dest"
)

func TestLocalRetainCountFiles(t *testing.T) {
	config := domain.ConfigDestLocal{
		Path: "/Users/jamestombleson",
	}

	c := dest.NewLocalRetainClient(config, "webdav", 1)
	_, err := c.CountFiles("retain")
	if err != nil {
		t.Error(err)
	}
}

func TestLocalRetainFindOldest(t *testing.T) {
	config := domain.ConfigDestLocal{
		Path: "/Users/jamestombleson",
	}

	c := dest.NewLocalRetainClient(config, "webdav", 1)
	oldest, err := c.FindOldestFile(".go")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("oldest: %v\n", oldest)
}

func TestLocalRetainCheck(t *testing.T) {
	config := domain.ConfigDestLocal{
		Path: "/Users/jamestombleson",
	}

	c := dest.NewLocalRetainClient(config, "webdav", 1)
	err := c.Check(".go")
	if err != nil {
		t.Error(err)
	}
}

func TestLocalMove(t *testing.T) {
	_, err := os.Create("fake.tar")
	if err != nil {
		t.Error(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	backupPath := filepath.Join(pwd, "fake.tar")

	c := dest.NewMoveClient("fake", backupPath, "test-container", pwd)
	err = c.Move(domain.RunDetails{})
	if err != nil {
		t.Error(err)
	}
	os.Remove("fake.tar")
	os.Remove(filepath.Join("test-container", "fake.tar"))
	os.Remove("test-container")
}
