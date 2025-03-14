package boardcmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NRF52840 操作
type NRF52840 struct {
	emprtprogram string
	scancmd      string
}

// Burn 烧写
func (driver *NRF52840) Burn(devport string, burnfile string, burncmd string) (string, error) {

	serialnum, err := GetJLinkSerialNumber(devport)
	if err != nil {
		log.Errorf("nRF52840 {%v} get jlink serial number error {%v}", devport, err)
	}

	cmdstr := fmt.Sprintf(burncmd, burnfile, serialnum)
	log.Infof("nRF52840 burn cmd {%v}", cmdstr)

	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Wait()
	stdoutStderr, err := cmd.CombinedOutput()
	log.Infof("nRF52840 burn result {%v}", string(stdoutStderr))
	if err != nil {
		log.Errorf("nRF52840 burn error {%v}", err)
		return fmt.Sprintf("%s\n\n%s", cmdstr, stdoutStderr), err
	}

	return fmt.Sprintf("%s\n\n%s", cmdstr, stdoutStderr), nil
}

// BurnEmptyProgram 烧写空程序
func (driver *NRF52840) BurnEmptyProgram(devport string, burncmd string) error {
	if output, err := driver.Burn(devport, driver.emprtprogram, burncmd); err != nil {
		err = fmt.Errorf("burn error {%v} {%v}", output, err)
		log.Error(err)
		return err
	}
	return nil
}

// SetEmptyProgram 设置空程序
func (driver *NRF52840) SetEmptyProgram(emptyprogram string) {
	driver.emprtprogram = emptyprogram
}

// Scan 扫描设备
func (driver *NRF52840) Scan() (map[string]bool, error) {

	cmd := exec.Command("sh", "-c", driver.scancmd)
	stdoutStderr, err := cmd.CombinedOutput()
	// log.Infof("nRF52840 scan output {%v}", string(stdoutStderr))
	if err != nil {
		log.Errorf("nRF52840 scan error {%v}", err)
		return nil, err
	}

	// 正则表达式匹配
	nummatch, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Errorf("regexp compile error {%v}", err)
		return nil, err
	}

	devmap := make(map[string]bool)
	template := "Serial Number:"
	for _, linestr := range strings.Split(string(stdoutStderr), "\n") {
		if strings.HasPrefix(linestr, template) {
			sernum := nummatch.FindString(linestr)
			devmap[fmt.Sprintf("/dev/%v-%v", "nRF52840", sernum)] = true
		}
	}

	return devmap, nil
}

// SetScanCmd 设置扫描命令
func (driver *NRF52840) SetScanCmd(scancmd string) {
	driver.scancmd = scancmd
}
