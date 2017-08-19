package utils

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

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

// IsDirMissingOrEmpty determines whether a directory is empty or missing
func IsDirMissingOrEmpty(path string) (bool, error) {
	dirExists, err := IsDirExist(path)
	if err != nil {
		return false, err
	}

	if !dirExists {
		return true, nil
	}

	dirEmpty, err := IsDirEmpty(path)
	if err != nil {
		return false, err
	}

	if dirEmpty {
		return true, nil
	}

	return false, nil
}

// IsDirEmpty determines whether a directory is empty
func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
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

// AppDataDir returns a default data directory for the databases
func AppDataDir() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if usr, err := user.Current(); err != nil {
			homeDir = usr.HomeDir
		}
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "deploy")
	default:
		return filepath.Join(homeDir, ".deploy")
	}
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

// FileExist determines whether a file exists
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
