package data

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	l        *log.Logger
	basePath string
	maxMem   int
}

func NewLocal(l *log.Logger, basePath string, maxFileSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Local{l, p, maxFileSize}, nil
}

func (local *Local) Save(filePath string, body io.Reader) error {
	local.l.Println("Storage interface's method save is running")

	fp := filepath.Join(local.basePath, filePath)

	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to create a directory: %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil {
		err := os.Remove(fp)
		if err != nil {
			return fmt.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Unable to get file info: %w", err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("Unable to create the file, error: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, body)
	if err != nil {
		return fmt.Errorf("Unable to write to the file: %w", err)
	}

	return nil
}
