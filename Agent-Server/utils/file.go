package utils

import(
	"os"
	"strings"
	"path/filepath"
	
)


func removeDir(walkDir string)error {
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