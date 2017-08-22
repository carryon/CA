package utils

import (
	"fmt"
	"os/exec"
)

func StartProcess(execFilePath, configFilePath string) error {
	startCmd := fmt.Sprintf("%s --config=%s  &", execFilePath, configFilePath)
	cmd := exec.Command("/bin/sh", "-c", startCmd)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func StopProcess(keyword string) error {
	killCmd := fmt.Sprintf("pkill -f  %s", keyword)
	cmd := exec.Command("/bin/sh", "-c", killCmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
