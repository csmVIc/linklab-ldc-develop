package device

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/parameter/response"
	"strings"

	log "github.com/sirupsen/logrus"
)

func userdevices(userid string) (*response.DeviceList, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		return nil, err
	}

	devClients, err := rdb.Keys(context.TODO(), "devices:active:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		return nil, err
	}

	boardnameDevMap := make(map[string][]response.DeviceStatus)
	for _, devClient := range devClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]
		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			return nil, err
		}
		for deviceid := range devices {
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if _, isOk := boardnameDevMap[boardname]; isOk == false {
				boardnameDevMap[boardname] = make([]response.DeviceStatus, 0)
			}
			boardnameDevMap[boardname] = append(boardnameDevMap[boardname], response.DeviceStatus{BoardName: boardname, DeviceID: deviceid, Busy: false, ClientID: clientid})
		}
	}

	useDevClients, err := rdb.Keys(context.TODO(), "devices:use:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		return nil, err
	}

	result := response.DeviceList{
		Devices: []response.DeviceStatus{},
	}
	for _, devClient := range useDevClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]
		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			return nil, err
		}
		for deviceid, devicestatusstr := range devices {
			// 解析开发板名
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if _, isOk := boardnameDevMap[boardname]; isOk == false {
				continue
			}
			// 解析设备状态
			deviceusestatus := value.DeviceUseStatus{}
			if err := json.Unmarshal([]byte(devicestatusstr), &deviceusestatus); err != nil {
				err = fmt.Errorf("json.Unmarshal error {%v}", err)
				log.Error(err)
				return nil, err
			}
			// 查找该设备的活跃状态
			if deviceusestatus.UserID == userid {
				i := -1
				for index, value := range boardnameDevMap[boardname] {
					if value.ClientID == clientid && value.DeviceID == deviceid {
						i = index
						result.Devices = append(result.Devices, response.DeviceStatus{
							BoardName: value.BoardName,
							DeviceID:  value.DeviceID,
							Busy:      true,
							ClientID:  value.ClientID,
						})
					}
				}
				if i == -1 {
					err := fmt.Errorf("device {%v:%v} is in use, but not in active error", clientid, deviceid)
					log.Error(err)
					return nil, err
				}
			}
		}
	}
	return &result, nil
}

func alldevcies(nboardname string) (*[]response.DeviceStatus, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		return nil, err
	}

	devClients, err := rdb.Keys(context.TODO(), "devices:active:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		return nil, err
	}

	boardnameDevMap := make(map[string][]response.DeviceStatus)
	for _, devClient := range devClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]
		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			return nil, err
		}
		for deviceid := range devices {
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if nboardname != "all" && nboardname != boardname {
				continue
			}
			if _, isOk := boardnameDevMap[boardname]; isOk == false {
				boardnameDevMap[boardname] = make([]response.DeviceStatus, 0)
			}
			boardnameDevMap[boardname] = append(boardnameDevMap[boardname], response.DeviceStatus{BoardName: boardname, DeviceID: deviceid, Busy: false, ClientID: clientid, Index: len(boardnameDevMap[boardname])})
		}
	}

	useDevClients, err := rdb.Keys(context.TODO(), "devices:use:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		return nil, err
	}

	for _, devClient := range useDevClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]
		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			return nil, err
		}
		for deviceid := range devices {
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if _, isOk := boardnameDevMap[boardname]; isOk == false {
				continue
			}
			i := -1
			for index, value := range boardnameDevMap[boardname] {
				if value.ClientID == clientid && value.DeviceID == deviceid {
					boardnameDevMap[boardname][index].Busy = true
					i = index
				}
			}
			if i == -1 {
				err := fmt.Errorf("device {%v:%v} is in use, but not in active error", clientid, deviceid)
				log.Error(err)
				return nil, err
			}
		}
	}

	var result []response.DeviceStatus
	for _, devices := range boardnameDevMap {
		result = append(result, devices...)
	}

	return &result, nil
}
