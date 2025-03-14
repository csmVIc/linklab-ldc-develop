package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ProjectDeploy(rootdirectory, kubeconfig, dockerconfig, projectname string) error {
	projectdirectory := filepath.Join(rootdirectory, projectname)
	log.Infof("project directory {%s}", projectdirectory)

	cmdstr := "build-noarch.sh"
	cmd := exec.Command("bash", cmdstr)
	cmd.Dir = projectdirectory
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CONFIG=%s", dockerconfig))
	outputbytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("command error {%v} output {%v}", err, string(outputbytes))
		return err
	}
	log.Infof("cmd {%s} output {%s}", cmdstr, string(outputbytes))

	cmdstr = "restart.sh"
	cmd = exec.Command("bash", cmdstr)
	cmd.Dir = projectdirectory
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBECONFIG=%s", kubeconfig))
	outputbytes, err = cmd.CombinedOutput()
	if err != nil {
		log.Errorf("command error {%v} output {%v}", err, string(outputbytes))
		return err
	}
	log.Infof("cmd {%s} output {%s}", cmdstr, string(outputbytes))

	return nil
}
