package setting

import (
	"bytes"
	"os/exec"
)

func RunShellCmd(name string, arg ...string) (error, string, string) {
	//非阻塞
	cmd := exec.Command(name, arg...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	return err, stderr.String(), out.String()

}
