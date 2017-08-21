package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func removeDir(walkDir string) error {
	var files []string
	err := filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return err
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.EqualFold(file, walkDir) {
			continue
		}
		os.RemoveAll(file)
	}

	return nil
}

// OpenFile opens or creates a file
// If the file already exists, open it . If it does not,
// It will create the file with mode 0644.
func OpenFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, err
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveFile(FilePath string) error {
	if err := os.Remove(FilePath); err != nil {
		return err
	}
	return nil
}

// OpenDir opens or creates a dir
// If the dir already exists, open it . If it does not,
// It will create the file with mode 0700.
func OpenDir(dir string) (string, error) {
	exists, err := IsDirExist(dir)
	if !exists {
		err = os.MkdirAll(dir, 0700)
	}
	return dir, err
}

// IsDirExist determines whether a directory exists
func IsDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExist determines whether a file exists
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
