package device

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/parameter/request"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
)

// idledevices 返回可用的空闲设备
// map[boardname]...DeviceIndex
func idledevices() (*map[string]*list.List, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return nil, err
	}

	devClients, err := rdb.Keys(context.TODO(), "devices:active:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%s}", err)
		log.Error(err)
		return nil, err
	}

	clientIDToTenantIDMap := make(map[string]map[int]bool)
	// boardnameIdleDevMap := make(map[string][]DeviceIndex)
	boardnameIdleDevMapList := make(map[string]*list.List)

	for _, devClient := range devClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]
		// 查询客户端所属的租户ID
		clientTenantID, err := cache.Cdriver.GetClientTenantID(clientid)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.GetClientTenantID {%v} error {%s}", clientid, err)
			log.Error(err)
			return nil, err
		}

		clientIDToTenantIDMap[clientid] = clientTenantID

		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			// 获取该客户端的活跃设备失败,直接跳过
			err = fmt.Errorf("redis hgetall error {%s}", err)
			log.Error(err)
			// return nil, err
			continue
		}
		for deviceid := range devices {
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if _, isOk := boardnameIdleDevMapList[boardname]; isOk == false {
				// boardnameIdleDevMap[boardname] = make([]DeviceIndex, 0)
				boardnameIdleDevMapList[boardname] = list.New()
			}
			// boardnameIdleDevMap[boardname] = append(boardnameIdleDevMap[boardname], DeviceIndex{
			// 	ClientID: clientid,
			// 	DeviceID: deviceid,
			// })
			if len(clientTenantID) == 1 {
				boardnameIdleDevMapList[boardname].PushFront(&DeviceIndex{
					ClientID: clientid,
					DeviceID: deviceid,
					Idle:     true,
				})
			} else {
				boardnameIdleDevMapList[boardname].PushBack(&DeviceIndex{
					ClientID: clientid,
					DeviceID: deviceid,
					Idle:     true,
				})
			}
		}
	}

	// boardnameIdleDevMap := make(map[string][]DeviceIndex)
	// for boardname := range boardnameIdleDevMapList {
	// 	boardnameIdleDevMap[boardname] = make([]DeviceIndex, 0)
	// 	for boardnameIdleDevMapList[boardname].Len() > 0 {
	// 		elem := boardnameIdleDevMapList[boardname].Front()
	// 		boardnameIdleDevMap[boardname] = append(boardnameIdleDevMap[boardname], *(elem.Value.(*DeviceIndex)))
	// 		boardnameIdleDevMapList[boardname].Remove(elem)
	// 	}
	// }

	useDevClients, err := rdb.Keys(context.TODO(), "devices:use:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%s}", err)
		log.Error(err)
		return nil, err
	}

	for _, devClient := range useDevClients {
		clientid := devClient[strings.LastIndex(devClient, ":")+1:]

		devices, err := rdb.HGetAll(context.TODO(), devClient).Result()
		if err != nil {
			// 获取该客户端的占用设备失败,直接跳过
			err = fmt.Errorf("redis hgetall error {%s}", err)
			log.Error(err)
			// return nil, err
			continue
		}
		for deviceid := range devices {
			boardname := deviceid[strings.LastIndex(deviceid, "/")+1 : strings.LastIndex(deviceid, "-")]
			if _, isOk := boardnameIdleDevMapList[boardname]; isOk == false {
				err := fmt.Errorf("boardname {%v} device {%v:%v} is in use, but not in active error", boardname, clientid, deviceid)
				log.Error(err)
				return nil, err
			}
			// i := -1
			// for index, value := range boardnameIdleDevMapList[boardname] {
			// 	if value.ClientID == clientid && value.DeviceID == deviceid {
			// 		i = index
			// 	}
			// }
			var elem *list.Element = nil
			for ptr := boardnameIdleDevMapList[boardname].Front(); ptr != nil; ptr = ptr.Next() {
				if ptr.Value.(*DeviceIndex).ClientID == clientid && ptr.Value.(*DeviceIndex).DeviceID == deviceid {
					elem = ptr
					break
				}
			}

			if elem == nil {
				err := fmt.Errorf("device {%v:%v} is in use, but not in active error", clientid, deviceid)
				log.Error(err)
				return nil, err
			}

			boardnameIdleDevMapList[boardname].Remove(elem)
			// boardnameIdleDevMap[boardname] = remove(boardnameIdleDevMap[boardname], i)

			if boardnameIdleDevMapList[boardname].Len() < 1 {
				delete(boardnameIdleDevMapList, boardname)
			}
		}
	}

	for _, blist := range boardnameIdleDevMapList {

		tmpList := []*DeviceIndex{}
		singleTenantOccupy := 0
		for ptr := blist.Front(); ptr != nil; ptr = ptr.Next() {
			value := ptr.Value.(*DeviceIndex)
			tmpList = append(tmpList, value)
			if len(clientIDToTenantIDMap[value.ClientID]) == 1 {
				singleTenantOccupy++
			}
		}

		if singleTenantOccupy > 0 {
			rand.Shuffle(len(tmpList[0:singleTenantOccupy]), func(i, j int) {
				tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
			})
		}

		rand.Shuffle(len(tmpList[singleTenantOccupy:]), func(i, j int) {
			tmpList[i+singleTenantOccupy], tmpList[j+singleTenantOccupy] = tmpList[j+singleTenantOccupy], tmpList[i+singleTenantOccupy]
		})

		blist.Init()

		for _, elem := range tmpList {
			blist.PushBack(elem)
		}
	}

	return &boardnameIdleDevMapList, err
}

// idlegroup 返回可用的空闲绑定组
func idlegroup(grouptype string) (int, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return -1, err
	}

	bindGroupMap, err := rdb.HGetAll(context.TODO(), fmt.Sprintf("bind:group:type:%v", grouptype)).Result()
	if err != nil {
		err = fmt.Errorf("redis hgetall error {%s}", err)
		log.Error(err)
		return -1, err
	}

	availGroups := []*GroupInfo_{}
	for bindgroupid := range bindGroupMap {
		deviceGroup := request.DeviceGroup{}
		if err := json.Unmarshal([]byte(bindGroupMap[bindgroupid]), &deviceGroup); err != nil {
			err = fmt.Errorf("json.Unmarshal error {%s}", err)
			log.Error(err)
			return -1, err
		}

		elem := GroupInfo_{
			ID:      bindgroupid,
			Devices: []DeviceIndex{},
		}

		for _, devinfo := range deviceGroup.Devices {
			// 检查设备是否活跃
			if isOk, err := rdb.HExists(context.TODO(), fmt.Sprintf("devices:active:%v", devinfo.ClientID), devinfo.DeviceID).Result(); err != nil {
				err = fmt.Errorf("rdb.HExists error {%s}", err)
				log.Error(err)
				return -1, err
			} else if isOk == false {
				break
			}

			// 检查设备是否空闲
			if isOk, err := rdb.HExists(context.TODO(), fmt.Sprintf("devices:use:%v", devinfo.ClientID), devinfo.DeviceID).Result(); err != nil {
				err = fmt.Errorf("rdb.HExists error {%s}", err)
				log.Error(err)
				return -1, err
			} else if isOk == true {
				break
			}

			elem.Devices = append(elem.Devices, DeviceIndex{
				ClientID: devinfo.ClientID,
				DeviceID: devinfo.DeviceID,
			})
		}

		// 验证通过
		if len(elem.Devices) == len(deviceGroup.Devices) {
			availGroups = append(availGroups, &elem)
		}
	}

	// if len(availGroups) < 1 {
	// 	return nil, errors.New("available groups length < 1 error")
	// }

	return len(availGroups), nil
}
