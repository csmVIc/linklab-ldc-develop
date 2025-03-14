package calltest

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func (cd *Driver) logmonitor(handle *websocket.Conn, logMap *map[string]interface{}, endchan chan<- error) {

	logbytescount := 0
	logrecvcount := 0
	avglogrecvresp := time.Duration(0)
	avgdevicelogrecvresp := time.Duration(0)
	avgnatslogrecvresp := time.Duration(0)
	avguserlogrecvresp := time.Duration(0)
	logs := []string{}
	for {
		taskdata := msg.TaskData{}
		resp := msg.UserMsg{
			Data: &taskdata,
		}

		err := handle.ReadJSON(&resp)
		if err != nil {
			err := fmt.Errorf("handle.ReadJSON error {%v}", err)
			log.Error(err)
			(*logMap)["websocketbroke"] = true
			endchan <- nil
			return
		}

		if resp.Code != 0 {
			err := fmt.Errorf("resp.Code {%v} != 0 error, {%v}", resp.Code, resp)
			log.Error(err)
			endchan <- err
		}

		switch resp.Type {
		case msg.TaskMsg:
			fmt.Println(taskdata)
			// taskdata := resp.Data.(msg.TaskData)
			switch taskdata.Type {
			case msg.TaskAllocateMsg:
				log.Infof("TaskAllocateMsg {%v}", taskdata)
				(*logMap)["deviceallocate"] = time.Now().UnixNano()
			case msg.TaskEndRunMsg:
				log.Infof("TaskEndRunMsg {%v}", taskdata)
				(*logMap)["logs"] = logs
				(*logMap)["websocketbroke"] = false
				(*logMap)["enddevicerun"] = time.Now().UnixNano()
				(*logMap)["logbytescount"] = logbytescount
				(*logMap)["logrecvcount"] = logrecvcount
				if logrecvcount == 0 {
					err := fmt.Errorf("logrecvcount == 0 error {%v}", err)
					log.Error(err)
					(*logMap)["websocketbroke"] = true
					endchan <- nil
					return
				}
				(*logMap)["avglogrecvresp"] = int(avglogrecvresp/time.Millisecond) / (logrecvcount)
				(*logMap)["avgdevicelogrecvresp"] = int(avgdevicelogrecvresp/time.Millisecond) / (logrecvcount)
				(*logMap)["avgnatslogrecvresp"] = int(avgnatslogrecvresp/time.Millisecond) / (logrecvcount)
				(*logMap)["avguserlogrecvresp"] = int(avguserlogrecvresp/time.Millisecond) / (logrecvcount)
				endchan <- nil
				return
			case msg.TaskBurnMsg:
				log.Infof("TaskBurnMsg {%v}", taskdata)
				(*logMap)["enddeviceburn"] = time.Now().UnixNano()
			case msg.TaskLogMsg:
				log.Infof("TaskLogMsg {%v}", taskdata)
				logs = append(logs, taskdata.Msg)
				logbytescount += len(taskdata.Msg)
				logrecvcount++
				dataMap := taskdata.Data.(map[string]interface{})
				log.Debugf("dataMap {%v}", dataMap)
				timestampint64, err := strconv.ParseInt(dataMap["ctimestamp"].(string), 10, 64)
				log.Debugln(timestampint64)
				if err != nil {
					err := fmt.Errorf("strconv.ParseInt error {%v}", err)
					log.Error(err)
					endchan <- err
				}
				creadtimestamp := time.Unix(0, timestampint64)
				avglogrecvresp += time.Now().Sub(creadtimestamp)

				timestampint64, err = strconv.ParseInt(dataMap["dtimestamp"].(string), 10, 64)
				log.Debugln(timestampint64)
				if err != nil {
					err := fmt.Errorf("strconv.ParseInt error {%v}", err)
					log.Error(err)
					endchan <- err
				}
				dreadtimestamp := time.Unix(0, timestampint64)
				avgdevicelogrecvresp += dreadtimestamp.Sub(creadtimestamp)

				ureadtimestamp := time.Unix(0, resp.TimeStamp)
				avgnatslogrecvresp += ureadtimestamp.Sub(dreadtimestamp)
				avguserlogrecvresp += time.Now().Sub(ureadtimestamp)

			default:
				err := fmt.Errorf("unknow TaskDataType {%v} {%v}", resp.Type, resp)
				log.Error(err)
				endchan <- err
				return
			}
		default:
			err := fmt.Errorf("unknow UserMsgType {%v} {%v}", resp.Type, resp)
			log.Error(err)
			endchan <- err
			return
		}
	}
}
