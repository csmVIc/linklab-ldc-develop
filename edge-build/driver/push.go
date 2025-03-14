package driver

import (
	"fmt"
	"os/exec"
)

func (dd *Driver) push() error {

	// 镜像推送
	cmdstr := fmt.Sprintf("docker push %v", dd.info.Build.ImageName)
	cmd := exec.Command("sh", "-c", cmdstr)
	stdoutByte, err := cmd.CombinedOutput()
	fmt.Printf("%s", string(stdoutByte))
	return err
}
