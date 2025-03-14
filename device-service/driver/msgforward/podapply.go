package msgforward

// func (md *Driver) podapplyhandler(subchansize int, msgtopic string, mqtttopic string) {

// 	defer func() {
// 		md.exitsignal = true
// 	}()

// 	natsconn, err := messenger.Mdriver.GetNatsConn()
// 	if err != nil {
// 		log.Errorf("get nats conn error {%v}", err)
// 		return
// 	}

// 	submsgchan := make(chan *nats.Msg, subchansize)
// 	defer func() {
// 		close(submsgchan)
// 	}()
// 	sub, err := natsconn.QueueSubscribe(msgtopic, "device-service", func(msg *nats.Msg) {
// 		submsgchan <- msg
// 		log.Infof("sub topic {%v} msg {%v}", msg.Subject, string(msg.Data))
// 	})
// 	defer func() {
// 		if err := sub.Unsubscribe(); err != nil {
// 			log.Errorf("unsubscribe messenger topic {%v} error {%v}", sub.Subject, err)
// 		}
// 	}()

// 	for {
// 		select {
// 		case nmsg := <-submsgchan:
// 			username, err := messenger.GetIDFromTopic(nmsg.Subject, "clients")
// 			if err != nil {
// 				log.Errorf("messenger.GetIDFromTopic error {%v}", err)
// 				return
// 			}
// 			clientid, err := auth.GetClientIDByUserName(username)
// 			if err != nil {
// 				log.Errorf("auth.GetClientIDByUserName error {%v}", err)
// 				return
// 			}

// 			// 转发至设备管理客户端
// 			err = mqttclient.MDriver.PubMsgByte(fmt.Sprintf(mqtttopic, username, clientid), 2, nmsg.Data)
// 			if err != nil {
// 				log.Errorf("mqttclient.MDriver.PubMsgByte error {%v}", err)
// 				return
// 			}
// 			log.Debugf("mqtt pub msg {%v} {%v}", fmt.Sprintf(mqtttopic, username, clientid), string(nmsg.Data))

// 			// 记录Pod部署日志
// 			// TODO

// 		case <-time.After(time.Second):
// 			if messenger.Mdriver.GetClosed() == true {
// 				log.Errorf("messenger.Mdriver.GetClosed true error")
// 				return
// 			}
// 		}
// 	}
// }
