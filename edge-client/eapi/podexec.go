package eapi

import (
	"fmt"
	"io"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/wsconf"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func podexec(c *gin.Context) {

	// 参数验证
	p := request.EdgeClientPodExec{}
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查Pod是否存在
	if pod := edgenode.EDriver.GetPod(p.Namespace, p.Pod); pod == nil {
		err := fmt.Errorf("pod {%v} not exist error", p.Pod)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	} else if len(p.Container) > 0 {
		isOk := false
		for _, container := range pod.Containers {
			if container.Name == p.Container {
				isOk = true
				break
			}
		}
		if isOk == false {
			err := fmt.Errorf("container {%v} not exist error", p.Container)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 建立回复
	protocols := websocket.Subprotocols(c.Request)
	var httpHeader http.Header = nil
	if len(protocols) > 0 {
		log.Debugf("websocket subprotocols {%v}", protocols)
		httpHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocols[0]},
		}
	}

	// 长连接
	wshander, err := wsconf.UpgraderGlobal.Upgrade(c.Writer, c.Request, httpHeader)
	if err != nil {
		err := fmt.Errorf("websocket upgrade error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer wshander.Close()

	// 准备管道
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	defer func() {
		stdinReader.Close()
		stdinWriter.Close()
		stdoutReader.Close()
		stdoutWriter.Close()
	}()

	// 退出通道
	readExitChan := make(chan error, 1)
	writeExitChan := make(chan error, 1)
	execExitChan := make(chan error, 1)

	// 读转发
	go func() {
		buf := make([]byte, 1024)
		for {
			len, err := stdoutReader.Read(buf)
			if err == io.EOF {
				continue
			}
			if err != nil {
				err = fmt.Errorf("stdoutReader.Read error {%v}", err)
				log.Error(err)
				readExitChan <- err
				return
			}
			if len > 0 {
				log.Debugf("stdoutReader.Read {%v}", string(buf[0:len]))
				if err := wshander.WriteJSON(&msg.EdgeClientPodExec{
					Type: msg.NormalPodExec,
					Msg:  string(buf[0:len]),
				}); err != nil {
					err = fmt.Errorf("wshander.WriteJSON error {%v}", err)
					log.Error(err)
					readExitChan <- err
					return
				}
			}
		}
	}()

	// 写转发
	go func() {
		for {
			execMsg := &msg.EdgeClientPodExec{}
			if err := wshander.ReadJSON(execMsg); err != nil {
				err = fmt.Errorf("wshander.ReadJSON error {%v}", err)
				log.Error(err)
				writeExitChan <- err
				return
			}
			if execMsg.Type == msg.ErrorPodExec {
				err = fmt.Errorf("wshander.ReadJSON error {%v}", execMsg.Msg)
				log.Error(err)
				writeExitChan <- err
				return
			}
			if _, err := stdinWriter.Write([]byte(execMsg.Msg)); err != nil {
				err = fmt.Errorf("stdinWriter.Write error {%v}", execMsg.Msg)
				log.Error(err)
				writeExitChan <- err
				return
			}
			log.Debugf("stdinWriter.Write {%v}", execMsg.Msg)
		}
	}()

	// Pod执行
	go func() {
		if err := edgenode.EDriver.PodExec(p.Namespace, p.Pod, p.Container, p.Commands, stdinReader, stdoutWriter); err != nil {
			// err = fmt.Errorf("edgenode.EDriver.PodExec error {%v}", err)
			log.Error(err)
			execExitChan <- err
			return
		}
		execExitChan <- nil
	}()

	// 监控
	for EXIST := false; !EXIST; {
		select {
		case err := <-readExitChan:
			log.Errorf("readExistChan error {%v}", err)
			EXIST = true
			break
		case err := <-writeExitChan:
			log.Errorf("writeExistChan error {%v}", err)
			if wserr := wshander.WriteJSON(&msg.EdgeClientPodExec{
				Type: msg.ErrorPodExec,
				Msg:  err.Error(),
			}); wserr != nil {
				log.Errorf("wshander.WriteJSON error {%v}", wserr)
			}
			EXIST = true
			break
		case err := <-execExitChan:
			if err != nil {
				log.Errorf("execExistChan error {%v}", err)
				if wserr := wshander.WriteJSON(&msg.EdgeClientPodExec{
					Type: msg.ErrorPodExec,
					Msg:  err.Error(),
				}); wserr != nil {
					log.Errorf("wshander.WriteJSON error {%v}", wserr)
				}
			}
			EXIST = true
			break
		}
	}
}
