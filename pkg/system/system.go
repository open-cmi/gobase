package system

import (
	"io"
	"os"
	"strings"
)

// ChangePasswd change os user passwd
func ChangeSystemPasswd(user string, encryptPasswd string) error {
	file, err := os.OpenFile("/etc/shadow", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var newLines []string = []string{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, user) {
			arr := strings.Split(line, ":")
			arr[1] = encryptPasswd
			newLine := strings.Join(arr, ":")
			newLines = append(newLines, newLine)
		} else {
			newLines = append(newLines, line)
		}
	}

	file.Seek(0, 0)
	for _, line := range newLines {
		file.WriteString(line + "\n")
	}
	file.Close()
	return nil
}
