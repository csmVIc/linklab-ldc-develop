package driver

import (
	"os/exec"
)

func (dd *Driver) clean() error {

	// 镜像清空
	cmdstr := "docker images -a | awk '{print $3}' | xargs docker rmi"
	cmd := exec.Command("sh", "-c", cmdstr)
	_, err := cmd.CombinedOutput()
	// fmt.Printf("%s\n", string(stdoutByte))
	return err
}
