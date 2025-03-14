package mqttclient

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Driver mqtt客户端
type Driver struct {
	info             *MInfo
	client           *mqtt.Client
	topicToSub       *map[string]SubTopic
	subInitErrSignal int32
	disconnCount     int32
}

// Driver mqtt客户端全局实例
var (
	MDriver *Driver
)

// Ping 检查mqtt客户端是否连接
func (md *Driver) ping() error {
	if md.client == nil {
		return errors.New("mqtt client nil error")
	}
	if token := (*md.client).Connect(); token.WaitTimeout(time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// 初始化mqtt客户端连接
func (md *Driver) init() error {
	// 根据当前进程的hostname和pid组合生成clientid
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("init os.Hostname error {%v}", err)
		return err
	}
	pid := os.Getpid()
	clientid := fmt.Sprintf("%s-%v", hostname, pid)

	// 如果是设备管理客户端,并且运行在云端,则替换用户名
	if md.info.Client.IsCloud {
		podname := os.Getenv("POD_NAME")
		bindex := strings.LastIndex(podname, "-")
		if bindex < 0 {
			err := fmt.Errorf("strings.Index os.Getenv(\"POD_NAME\") {%v} < {%v} 0", podname, bindex)
			log.Error(err)
			return err
		}
		md.info.Client.UserName = md.info.Client.UserName + podname[bindex:]
	}

	// mqtt客户端参数初始化
	opts := mqtt.NewClientOptions().
		AddBroker(md.info.Client.URL).
		SetAutoReconnect(true).
		SetConnectTimeout(time.Second).
		SetMaxReconnectInterval(time.Second).
		SetKeepAlive(time.Second * 60).
		SetPingTimeout(time.Second * 20).
		SetUsername(md.info.Client.UserName).
		SetPassword(md.info.Client.PassWord).
		SetProtocolVersion(4).
		SetClientID(clientid).
		SetDefaultPublishHandler(md.msghandler).
		SetConnectionLostHandler(md.connlosthandler).
		SetOnConnectHandler(md.onconnhandler)

	// mqtt客户端创建
	client := mqtt.NewClient(opts)
	md.client = &client
	if err := md.ping(); err != nil {
		log.Errorf("mqtt ping error {%v}", err)
		return err
	}

	return nil
}

// New 创建mqtt客户端实例
func New(i *MInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("info is nil error")
	}
	md := &Driver{info: i, client: nil, topicToSub: nil, subInitErrSignal: 0, disconnCount: 0}
	if err := md.init(); err != nil {
		log.Errorf("mqtt client init error {%v}", err)
		return nil, err
	}
	return md, nil
}
