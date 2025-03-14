package topichandler

// func (td *Driver) podapplysub(client mqtt.Client, mqttmsg mqtt.Message) {

// 	var err error = nil
// 	defer func() {
// 		if err != nil {
// 			*td.errchan <- err
// 		}
// 	}()

// 	var p msg.PodApply
// 	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
// 		log.Errorf("json Unmarshal error {%v}", err)
// 		return
// 	}

// 	select {
// 	case (*td.podapplychan) <- &p:
// 		log.Debugf("pod apply msg into podapplychan {%v}", p)
// 		return
// 	case <-time.After(time.Second * time.Duration(td.info.PodApply.ChanTimeOut)):
// 		err = fmt.Errorf("pod apply msg into podapplychan timeout error {%vs}", td.info.DeviceBurn.ChanTimeOut)
// 		log.Error(err)
// 		return
// 	}
// }

// func (td *Driver) podapplyresultresfuse(client mqtt.Client, mqttmsg mqtt.Message) {
// 	var err error = nil
// 	defer func() {
// 		if err != nil {
// 			*td.errchan <- err
// 		}
// 	}()

// 	var p msg.ReplyMsg
// 	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
// 		log.Errorf("json Unmarshal error {%v}", err)
// 		return
// 	}

// 	err = fmt.Errorf("mqtt pod apply result refuse {%v}", p.Msg)
// 	log.Error(err)
// }

// // PubPodApplyResult 发布Pod部署结果
// func (td *Driver) PubPodApplyResult(groupid string, success bool, msgstr string, beginapplytime int64, endapplytime int64) error {

// 	parameter := msg.PodApplyResult{
// 		GroupID:        groupid,
// 		Success:        success,
// 		Msg:            msgstr,
// 		BeginApplyTime: beginapplytime,
// 		EndApplyTime:   endapplytime,
// 	}

// 	err := mqttclient.MDriver.PubMsg((*td.topicMap)["podapply"].Pub, 2, parameter)
// 	if err != nil {
// 		log.Errorf("mqtt client pub msg error {%v}", err)
// 		return err
// 	}

// 	return nil
// }
