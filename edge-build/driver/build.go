package driver

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

func (dd *Driver) build() error {

	// 镜像打包推送
	cmdstr := fmt.Sprintf("docker build -t %v --file=Dockerfile .", dd.info.Build.ImageName)
	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Dir = "/app/workspace"
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("cmd.StdoutPipe error {%v}", err)
		log.Error(err)
		return err
	}
	linereader := bufio.NewScanner(stdout)

	time.Sleep(time.Second)
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("cmd.Start error {%v}", err)
		log.Error(err)
		return err
	}

	for linereader.Scan() {
		line := linereader.Text()
		fmt.Println(line)
	}

	if err := cmd.Wait(); err != nil {
		// err = fmt.Errorf("cmd.Wait error {%v}", err)
		// log.Error(err)
		return err
	}

	return nil
}
