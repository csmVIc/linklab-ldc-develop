package iotnode

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (id *Driver) setlogerr(devport string, err error, devstatus *DeviceStatus, taskchan chan<- *TaskRunEnd) error {
	// 如果打开串口错误,则设置错误状态
	tmpstatus := devstatus.BusyStatus
	devstatus.BusyStatus = LogError
	// 如果任务运行未结束,则需要发送结束信息
	if tmpstatus == Running {
		select {
		case taskchan <- &TaskRunEnd{BurnInfo: devstatus.BurnInfo, TimeOut: 2, EndTime: time.Now()}:
		case <-time.After(time.Duration(id.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
			log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devstatus.BurnInfo.DeviceID, id.info.DeviceLog.TaskTimeOutMill)
		}
	}
	// 设置错误信息
	log.Error(err)
	// 传递设备错误信息
	select {
	case id.devices.deviceErrorChan <- &DeviceError{DeviceID: devport, Err: err}:
	case <-time.After(time.Duration(id.info.DeviceError.TimeOut) * time.Second):
		err := fmt.Errorf("device status into err chan timeout {%v}", id.info.DeviceError.TimeOut)
		log.Error(err)
		return err
	}
	return nil
}

// ReadDeviceLog 读取设备日志
func (id *Driver) ReadDeviceLog(devport string, taskchan chan<- *TaskRunEnd) (<-chan *LogMsg, error) {

	var err error
	logchan := make(chan *LogMsg, id.info.DeviceLog.LogChanSize)
	defer func() {
		if err != nil {
			close(logchan)
		}
	}()

	board, err := id.GetBoardFromDevPort(devport)
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

		logbuffer := ""
		loglastsendtime := time.Now()
		var devicestatus *DeviceStatus = nil

		for {
			// 检测设备是否存在
			value, isExist := id.devices.devicesMap.Load(devport)
			if isExist == false {
				err = fmt.Errorf("devport {%v} not exist error", devport)
				log.Error(err)
				// 如果最近的状态存在,并且为运行状态,发送任务结束运行信号
				if devicestatus != nil && devicestatus.BusyStatus == Running {
					select {
					case taskchan <- &TaskRunEnd{BurnInfo: devicestatus.BurnInfo, TimeOut: 2, EndTime: time.Now()}:
					case <-time.After(time.Duration(id.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
						log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devport, id.info.DeviceLog.TaskTimeOutMill)
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
					// 烧写空程序
					if err := id.boardCmdMap[board].BurnEmptyProgram(devport, id.info.Commands[board].Burn, devicestatus.SerialPort); err != nil {
						err := id.setlogerr(devport, fmt.Errorf("devport {%v} burn empty program {%v} error {%v}", devport, id.info.Commands[board].NetworkCmd, err), devicestatus, taskchan)
						if err != nil {
							log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
						}
					}

					// 发送结束状态
					devicestatus.BusyStatus = IdleState
					tmpburninfo := devicestatus.BurnInfo
					devicestatus.BurnInfo = nil
					select {
					case taskchan <- &TaskRunEnd{BurnInfo: tmpburninfo, TimeOut: 1, EndTime: time.Now()}:
					case <-time.After(time.Duration(id.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
						log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devport, id.info.DeviceLog.TaskTimeOutMill)
					}
				}
			}

			// 空闲状态,重置日志缓存
			if devicestatus.BusyStatus == IdleState {
				logbuffer = ""
			}

			// 判断设备串口是否是开启状态,如果没有那就开启设备串口
			if err := id.boardCmdMap[board].ReOpen(devicestatus.SerialPort); err != nil {
				err := id.setlogerr(devport, fmt.Errorf("devport {%v} reopen error {%v}", devport, err), devicestatus, taskchan)
				if err != nil {
					log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
				}
				// 设备操作解锁
				devicestatus.Lock.Unlock()
				return
			}

			// 检查是否有命令需要写入
			for cmdreadexist := false; !cmdreadexist; {
				select {
				case devicecmd := <-(*devicestatus.CmdChan):
					// 检查设备是否为运行
					if devicestatus.BusyStatus != Running && devicestatus.BusyStatus != Burning {
						log.Debugf("device status not running or burning, so device cmd skip {%v}", devicecmd.Cmd)
						break
					}
					// 检查任务ID是否相同
					if devicestatus.BurnInfo.GroupID != devicecmd.GroupID || devicestatus.BurnInfo.TaskIndex != devicecmd.TaskIndex {
						log.Debugf("device status {%v:%v} not equal device cmd {%v:%v}, so skip", devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, devicecmd.GroupID, devicecmd.TaskIndex)
						break
					}
					// 命令写入
					if _, err := devicestatus.SerialPort.Write([]byte(devicecmd.Cmd)); err != nil {
						err := id.setlogerr(devport, fmt.Errorf("devport {%v} cmd write {%v} error {%v}", devport, devicecmd.Cmd, err), devicestatus, taskchan)
						if err != nil {
							log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
						}
						devicestatus.Lock.Unlock()
						return
					}
				default:
					cmdreadexist = true
				}
			}

			isready, err := devicestatus.SerialPort.ReadyToRead()
			if err != nil {
				err := id.setlogerr(devport, fmt.Errorf("devport {%v} ready to read error {%v}", devport, err), devicestatus, taskchan)
				if err != nil {
					log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
				}
				// 设备操作解锁
				devicestatus.Lock.Unlock()
				return
			}

			if isready < 1 {
				// 设备操作解锁
				devicestatus.Lock.Unlock()
				// 在没有日志可以读的情况下,睡眠一段时间减少CPU占用
				time.Sleep(time.Millisecond * time.Duration(id.info.DeviceLog.ReadSleepMill))
				// 更新日志读时间
				loglastsendtime = time.Now()
				continue
			}

			tmpbuffer := make([]byte, 1024)
			readlen, err := devicestatus.SerialPort.Read(tmpbuffer)
			if err != nil {
				err := id.setlogerr(devport, fmt.Errorf("devport {%v} read error {%v}", devport, err), devicestatus, taskchan)
				if err != nil {
					log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
				}
				// 设备操作解锁
				devicestatus.Lock.Unlock()
				return
			}

			// 设备运行状态才有读取日志的必要
			if devicestatus.BusyStatus != Running {
				devicestatus.Lock.Unlock()
				// 设备未在运行状态,可以睡眠一段时间减少CPU占用
				time.Sleep(time.Millisecond * time.Duration(id.info.DeviceLog.ReadSleepMill))
				// 更新日志读时间
				loglastsendtime = time.Now()
				continue
			}

			// 设备操作解锁
			devicestatus.Lock.Unlock()

			if readlen > 0 {
				// 添加到日志缓存
				logbuffer = logbuffer + string(tmpbuffer[:readlen])
				// 解析出日志缓存中的一行
				lines := strings.Split(logbuffer, "\n")
				for index := 0; index < len(lines)-1; index++ {

					// 网络初始化扫描
					if len(id.info.Commands[board].NetworkScan) > 0 {
						if strings.Index(lines[index], id.info.Commands[board].NetworkScan) != -1 {
							// 发送网络初始化命令
							if _, err := devicestatus.SerialPort.Write([]byte(id.info.Commands[board].NetworkCmd)); err != nil {
								err := id.setlogerr(devport, fmt.Errorf("devport {%v} set network {%v} error {%v}", devport, id.info.Commands[board].NetworkCmd, err), devicestatus, taskchan)
								if err != nil {
									log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
								}
								return
							}
						}
					}

					select {
					case logchan <- &LogMsg{
						Msg:       lines[index] + "\n",
						TimeStamp: time.Now().UnixNano(),
					}:
					case <-time.After(time.Millisecond * time.Duration(id.info.DeviceLog.LogTimeOutMill)):
						log.Errorf("devport {%v} log {%v} into logchan timeout {%vmills} error", devport, lines[index], id.info.DeviceLog.LogTimeOutMill)
						continue
					}
					loglastsendtime = time.Now()

					// 如果程序输出end,则设置为主动释放设备
					if devicestatus.BusyStatus == Running && devicestatus.BurnInfo != nil {
						if strings.Index(lines[index], "end") == 0 && board != "Haas100Python" {
							time.Sleep(time.Second)

							// 烧写空程序
							if err := id.boardCmdMap[board].BurnEmptyProgram(devport, id.info.Commands[board].Burn, devicestatus.SerialPort); err != nil {
								err := id.setlogerr(devport, fmt.Errorf("devport {%v} burn empty program {%v} error {%v}", devport, id.info.Commands[board].NetworkCmd, err), devicestatus, taskchan)
								if err != nil {
									log.Errorf("devport {%v} setlogerr error {%v}", devport, err)
								}
							}

							// 发送结束信息
							devicestatus.BusyStatus = IdleState
							tmpburninfo := devicestatus.BurnInfo
							devicestatus.BurnInfo = nil
							select {
							case taskchan <- &TaskRunEnd{BurnInfo: tmpburninfo, TimeOut: 0, EndTime: time.Now()}:
							case <-time.After(time.Duration(id.info.DeviceLog.TaskTimeOutMill) * time.Millisecond):
								log.Errorf("devport {%v} task into taskchan timeout {%vmills} error", devport, id.info.DeviceLog.TaskTimeOutMill)
							}
						}
					}
				}
				// 更新日志缓存
				logbuffer = lines[len(lines)-1]

				// 如果长时间日志未发送,我们就不等待整行日志,而是强制进行日志发送
				if len(logbuffer) > 0 && loglastsendtime.Add(time.Duration(id.info.DeviceLog.LogSendTimeOutMill)*time.Millisecond).After(time.Now()) == false {
					select {
					case logchan <- &LogMsg{
						Msg:       logbuffer,
						TimeStamp: time.Now().UnixNano(),
					}:
					case <-time.After(time.Millisecond * time.Duration(id.info.DeviceLog.LogTimeOutMill)):
						log.Errorf("devport {%v} log {%v} into logchan timeout {%vmills} error", devport, logbuffer, id.info.DeviceLog.LogTimeOutMill)
					}
					logbuffer = ""
					loglastsendtime = time.Now()
				}
			}
		}
	}()

	return logchan, nil
}
