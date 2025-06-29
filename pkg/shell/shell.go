package shell

import "os/exec"

func Execute(cmd string) error {
	args := []string{"-c", cmd}
	command := exec.Command("bash", args...)
	return command.Run()
}

func ExecuteOutput(cmd string) ([]byte, error) {
	args := []string{"-c", cmd}
	command := exec.Command("bash", args...)
	return command.Output()
}
