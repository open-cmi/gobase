package systemctl

import "os/exec"

func StartService(service string) error {
	args := []string{"enable", service}
	cmd := exec.Command("systemctl", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	args = []string{"start", service}
	cmd = exec.Command("systemctl", args...)
	return cmd.Run()
}

func StopService(service string) error {
	args := []string{"stop", service}
	cmd := exec.Command("systemctl", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	args = []string{"disable", service}
	cmd = exec.Command("systemctl", args...)
	return cmd.Run()
}

func StatusService(service string) bool {
	args := []string{"is-active", service}
	cmd := exec.Command("systemctl", args...)
	err := cmd.Run()
	return err == nil
}

func RestartService(service string) error {
	args := []string{"restart", service}
	cmd := exec.Command("systemctl", args...)
	err := cmd.Run()
	return err
}
