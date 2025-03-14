package virtualnode

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ReadDeviceLog 读取设备日志
func (vd *Driver) ReadDeviceLog(devport string, taskchan chan<- *TaskRunEnd) (<-chan *LogMsg, error) {

	var err error
	logchan := make(chan *LogMsg, vd.info.DeviceLog.LogChanSize)
	defer func() {
		if err != nil {
			close(logchan)
		}
	}()

	board, err := vd.GetBoardFromDevPort(devport)
	if err != nil {
		err = fmt.Errorf("get board from devport {%v} error {%v}", devport, err)
		log.Error(err)
		return logchan, err
	}

	go func() {
		// 退出时关闭log通道
		defer func() {
			close(logchan)
		}()

		var devicestatus *DeviceStatus = nil

		// 日志编号
		lognum := int64(0)

		for {
			// 检测设备是否存在
			value, isExist := vd.devices.devicesMap.Load(devport)
			if isExist == false {
				err = fmt.Errorf("devport {%v} not exist error", devport)
				log.Error(err)
				// 如果最近的状态存在,并且为运行状态,发送任务结束运行信号
				if devicestatus != nil && devicestatus.BusyStatus == Running {
					select {
					case taskchan <- &TaskRunEnd{BurnInfo: devicestatus.BurnInfo, TimeOut: 2, EndTime: time.Now()}:
					case <-time.After(time.Duration(vd.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
						log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devport, vd.info.DeviceLog.TaskTimeOutMill)
					}
				}
				return
			}
			devicestatus = value.(*DeviceStatus)

			// 设备操作加锁
			devicestatus.Lock.Lock()

			// 检查运行状态,如果运行超时,则设置为空闲状态
			if devicestatus.BusyStatus == Running && devicestatus.BurnInfo != nil {
				if devicestatus.BeginTime.Add(time.Duration(devicestatus.BurnInfo.RunTime)*time.Second).After(time.Now()) == false {
					devicestatus.BusyStatus = IdleState
					tmpburninfo := devicestatus.BurnInfo
					devicestatus.BurnInfo = nil
					select {
					case taskchan <- &TaskRunEnd{BurnInfo: tmpburninfo, TimeOut: 1, EndTime: time.Now()}:
					case <-time.After(time.Duration(vd.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
						log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devport, vd.info.DeviceLog.TaskTimeOutMill)
					}
				}
			}

			// 设备运行状态才有读取日志的必要
			if devicestatus.BusyStatus != Running {
				devicestatus.Lock.Unlock()
				time.Sleep(time.Millisecond * time.Duration(vd.info.DeviceLog.ReadSleepMill))
				continue
			}

			// 产生的字符串
			var sb strings.Builder
			sb.WriteString("%")
			sb.WriteString(strconv.Itoa(vd.info.Boards[board].LogBytes))
			sb.WriteString("d")
			lognum++
			tmpstr := fmt.Sprintf(sb.String(), lognum)

			// 设备操作解锁
			devicestatus.Lock.Unlock()

			// 日志发送
			if len(tmpstr) > 0 {
				select {
				case logchan <- &LogMsg{
					Msg:       tmpstr,
					TimeStamp: time.Now().UnixNano(),
				}:
				case <-time.After(time.Millisecond * time.Duration(vd.info.DeviceLog.LogTimeOutMill)):
					log.Errorf("devport {%v} log {%v} into logchan timeout {%vmills} error", devport, tmpstr, vd.info.DeviceLog.LogTimeOutMill)
					continue
				}
			}

			time.Sleep(time.Millisecond * time.Duration(vd.info.DeviceLog.ReadSleepMill))
		}
	}()

	return logchan, nil
}
