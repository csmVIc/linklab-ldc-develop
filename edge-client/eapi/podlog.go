package eapi

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/wsconf"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func podlog(c *gin.Context) {

	// 参数验证
	p := request.EdgeClientPodLog{}
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
    log.Debugf("收到Pod日志请求，命名空间: %s, pod名称: %s, 容器: %s", p.Namespace, p.Pod, p.Container)

	// 建立回复
	protocols := websocket.Subprotocols(c.Request)
	var httpHeader http.Header = nil
	if len(protocols) > 0 {
        log.Debugf("websocket子协议: {%v}", protocols)
		httpHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocols[0]},
		}
	}

	// 长连接
    log.Debugf("尝试升级连接为WebSocket")
	wshander, err := wsconf.UpgraderGlobal.Upgrade(c.Writer, c.Request, httpHeader)
	if err != nil {
        err := fmt.Errorf("websocket升级错误 {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer wshander.Close()
    log.Debugf("WebSocket连接建立成功")

	// 2. 设置更大的消息缓冲区
	log.Debugf("设置WebSocket超时时间和缓冲区大小")
	wshander.SetReadDeadline(time.Now().Add(time.Hour))    // 设置读取超时为1小时
	wshander.SetWriteDeadline(time.Now().Add(time.Hour))   // 设置写入超时为1小时
	wshander.SetReadLimit(1024 * 1024)
	messageCount := 0  // 消息计数器
	// 3. 修改心跳机制
    log.Debugf("配置ping处理器")
    wshander.SetPingHandler(func(message string) error {
        log.Debugf("收到ping请求，响应pong")
        err := wshander.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second*10))
        if err == websocket.ErrCloseSent {
            log.Debugf("发送pong时连接已关闭")
            return nil
        }
		if err != nil {
            log.Errorf("发送pong失败: %v", err)
        }
        return err
    })
	
	// 通道准备
    log.Debugf("初始化通道，缓冲区大小: %d", einfo.PodLog.OutChanSize)
	outputchan := make(chan string, einfo.PodLog.OutChanSize)
	exitchan := make(chan bool, 1)
	endchan := make(chan error, 1)
	defer func() {
		log.Debugf("开始清理资源，当前消息计数: %d", messageCount)
		exitchan <- true
		time.Sleep(time.Second)
		// 检查channel状态
		log.Debugf("输出通道当前长度: %d", len(outputchan))
		close(outputchan)
		close(exitchan)
		close(endchan)
		log.Debugf("资源清理完成")
	}()

	// 读取日志
	go func() {
        log.Debugf("启动日志读取协程")
		defer func() {
            log.Debugln("日志读取协程退出")
		}()
		log.Debugf("开始获取Pod日志")
		err = edgenode.EDriver.GetPodLog(p.Namespace, p.Pod, p.Container, outputchan, exitchan)
		if err != nil {
			// err := fmt.Errorf("%v", err)
            log.Errorf("获取Pod日志错误: %v", err)
			endchan <- err
			return
		}else{
			log.Debugf("Pod日志获取完成，无错误")
		}
        log.Debugf("日志读取成功完成")
		endchan <- nil
		return
	}()

	defer func() {
        log.Debugln("主podlog函数退出")
	}()

	// 转发消息
    log.Debugf("开始主消息转发循环")
	for {
		select {
		case err := <-endchan:
			if err != nil {
                log.Errorf("错误通道收到错误: %v", err)
				resp := &response.EdgeClientPodLog{
					Type: response.ErrorPodLog,
					Msg:  err.Error(),
				}
				if err := wshander.WriteJSON(resp); err != nil {
                    log.Errorf("向WebSocket写入错误消息失败: %v", err)
					return
				}
			}
			log.Debugf("总共转发消息数: %d", messageCount)
			return
		case logline := <-outputchan:
            messageCount++
            // 每处理100条消息记录一次日志
			if messageCount%100 == 0 {
				// log.Debugf("已转发消息数: %d，最新消息长度: %d", messageCount, len(logline))
			}
			// 每次收到新的日志时更新超时时间
			wshander.SetReadDeadline(time.Now().Add(time.Hour))
			wshander.SetWriteDeadline(time.Now().Add(time.Hour))

			resp := &response.EdgeClientPodLog{
				Type: response.NormalPodLog,
				Msg:  logline,
			}
			if err := wshander.WriteJSON(resp); err != nil {
				log.Errorf("向WebSocket写入日志消息失败: %v，消息长度: %d", err, len(logline))
				return
			}
		case <-time.After(time.Second * 30):
			// log.Debugf("发送ping心跳 (当前消息计数: %d)", messageCount)
			// 更新超时时间
			wshander.SetReadDeadline(time.Now().Add(time.Hour))
			wshander.SetWriteDeadline(time.Now().Add(time.Hour))
			if err := wshander.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
                log.Errorf("发送ping消息失败: %v", err)
				return
			}
		}
	}
}
