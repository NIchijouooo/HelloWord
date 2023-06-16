package setting

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

func Exec_shell(s_cmd string) (string, bool) {
	cmd := exec.Command("/bin/bash", "-c", s_cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd StdoutPipe error:", err)
		return "", false
	}
	cmd.Start()

	var end_line string
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		end_line += line
	}

	cmd.Wait()
	return end_line, true
}
