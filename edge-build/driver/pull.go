package driver

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (dd *Driver) pull() error {

	// 镜像加载
	tmpfile, err := ioutil.ReadFile("/app/workspace/Dockerfile")
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile error {%v}", err)
		log.Error(err)
		return err
	}

	// 镜像名解析
	imagelist := []string{}
	for index := 0; index < len(tmpfile); {
		index = strings.Index(string(tmpfile[index:]), "FROM")
		if index < 0 {
			break
		}
		index += 5
		jndex := index
		for ; tmpfile[jndex] != '\r' && tmpfile[jndex] != '\n'; jndex++ {
		}
		imagename := string(tmpfile[index:jndex])
		imagelist = append(imagelist, imagename)
		index = jndex + 1
		log.Debugf("resolve image name: %v", imagename)
	}

	// 镜像拉取
	for _, imagename := range imagelist {
		log.Debugf("docker pull image: %v", imagename)

		cmdstr := fmt.Sprintf("docker pull %v", imagename)
		// 调用系统shell程序通常是/bin/sh
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
	}

	return nil
}
