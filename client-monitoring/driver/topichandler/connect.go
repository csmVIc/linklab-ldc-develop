package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) clientconnect(client mqtt.Client, mqttmsg mqtt.Message) {

	clientMap := map[string]interface{}{}
	if err := json.Unmarshal(mqttmsg.Payload(), &clientMap); err != nil {
		log.Errorf("json.Unmarshal error {%v}", err)
		return
	}

	username := clientMap["username"].(string)
	clientid := clientMap["clientid"].(string)
	log.Debugf("mqtt client conn username {%v} clientid {%v}", username, clientid)

	// 日志记录
	loginsuccess := false
	defer func() {
		tags := map[string]string{
			"username": username,
		}
		fields := map[string]interface{}{
			"clientid": clientid,
			"success":  loginsuccess,
			"podname":  os.Getenv("POD_NAME"),
			"nodename": os.Getenv("NODE_NAME"),
		}
		err := logger.Ldriver.WriteLog("clientconnect", tags, fields)
		if err != nil {
			err = fmt.Errorf("write log err {%v}", err)
			log.Error(err)
			return
		}
	}()

	// 查找缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return
	}

	// 检查是否已经登录
	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("clients:id:%s:token:*", username)).Result()
	if err == nil {
		if len(tkeys) == 1 {

			valuestr, err := rdb.Get(context.TODO(), tkeys[0]).Result()
			if err != nil {
				log.Errorf("{%v} redis get error {%v}", username, err)
				return
			}

			clientloginstatus := value.ClientLoginStatus{}
			if err := json.Unmarshal([]byte(valuestr), &clientloginstatus); err != nil {
				log.Errorf("{%v} json.Unmarshal error {%v}", username, err)
				return
			}

			// 已经登录检查clientid是否相同
			if clientloginstatus.ClientID != clientid {
				// 已经有客户端登录,并且clientid并不相同
				err = fmt.Errorf("{%v} already logged in, but {%v} != {%v} error", username, clientloginstatus.ClientID, clientid)
				log.Error(err)
				if err := mqttclient.MDriver.PubMsg(fmt.Sprintf(td.info.Login.TokenRefuseTopic, td.info.Login.TokenPubDelay, username, clientid), 0, msg.ReplyMsg{Code: -1, Msg: err.Error()}); err != nil {
					log.Errorf("{%s} publish token refuse error {%v}", username, err)
					return
				}
				return
			}

			// 如果相同则需要延长token的时间
			log.Infof("clients {%s} {%s} already logged in, token {%v}\n", username, clientid, tkeys[0])
			_, err = rdb.Expire(context.TODO(), tkeys[0], time.Second*time.Duration(td.info.Login.TTL)).Result()
			if err != nil {
				log.Errorf("{%s} redis expire error {%v}", username, err)
				return
			}

			// 解析token
			token, err := auth.GetTokenFromKey(tkeys[0])
			if err != nil {
				log.Errorf("{%s} get token from key error {%v}", username, err)
				return
			}

			// 返回token报文
			reply := msg.ReplyMsg{
				Code: 0,
				Msg:  "login success",
				Data: map[string]string{
					"token": token,
				},
			}
			if err := mqttclient.MDriver.PubMsg(fmt.Sprintf(td.info.Login.TokenPubTopic, td.info.Login.TokenPubDelay, username, clientid), 0, reply); err != nil {
				log.Errorf("{%s} publish token error {%v}", username, err)
				return
			}
			loginsuccess = true
			log.Infof("client {%s} {%s} login success", username, clientid)
			return
		} else if len(tkeys) > 1 {
			// 如果发现多个token,则系统可能出现问题
			log.Errorf("clients {%s} {%s} already logged in, has multiple token error {%v}\n", username, clientid, tkeys)
			return
		}
	}

	// 未登录
	cinfo := &table.Client{}
	filter := table.ClientFilter{UserName: username}
	if err := database.Mdriver.FindOneElem("clients", filter, cinfo); err != nil {
		log.Errorf("{%v} database.Mdriver.FindOneElem error {%v}", username, err)
		return
	}

	// 检查是否为超级用户
	if cinfo.IsSuperUser == true {
		log.Debugf("client {%v} is super user, so skip", cinfo.UserName)
		return
	}

	clientloginstatus := value.ClientLoginStatus{
		ClientID: clientid,
		TenantID: cinfo.TenantID,
	}
	valuebyte, err := json.Marshal(clientloginstatus)
	if err != nil {
		log.Errorf("{%v} json.Marshal error {%v}", username, err)
		return
	}

	// 设置缓存
	token := auth.CreateToken(cinfo.UserName, cinfo.Password)
	_, err = rdb.Set(context.TODO(), fmt.Sprintf("clients:id:%s:token:%s", cinfo.UserName, token), string(valuebyte), time.Second*time.Duration(td.info.Login.TTL)).Result()
	if err != nil {
		log.Errorf("{%v} redis set error {%v}", username, err)
		return
	}

	// 返回token报文
	reply := msg.ReplyMsg{
		Code: 0,
		Msg:  "login success",
		Data: map[string]string{
			"token": token,
		},
	}
	if err := mqttclient.MDriver.PubMsg(fmt.Sprintf(td.info.Login.TokenPubTopic, td.info.Login.TokenPubDelay, username, clientid), 0, reply); err != nil {
		log.Errorf("{%s} publish token error {%v}", username, err)
		return
	}

	loginsuccess = true
	log.Infof("client {%s} {%s} login success", username, clientid)
	return
}
