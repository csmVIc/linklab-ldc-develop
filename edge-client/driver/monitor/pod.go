package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) pod() {

	defer func() {
		// 如果该函数退出,则将结束运行标志位设置为true
		md.endrun = true
		log.Error("pod func set endrun")
	}()

	podlastupdatetime := time.Now()
	for md.endrun == false {
		// 检查Pod变化
		podsChange, err := edgenode.EDriver.GetPodsChange()
		if err != nil {
			err = fmt.Errorf("edgenode.EDriver.GetPodsChange error {%v}", err)
			log.Error(err)
			return
		}

		if len(podsChange.Delete) > 0 || len(podsChange.HeartBeat) > 0 {
			// Pod列表发生变化
			if err := topichandler.TDriver.PubPodUpdate(podsChange); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubPodUpdate error {%v}", err)
				log.Error(err)
				return
			}
		}

		// 定时上传Pod列表，即使Pod列表没有变化，也上传Pod列表
		if time.Now().After(podlastupdatetime) {
			pods := edgenode.EDriver.GetPods()
			// 更新Pod列表
			if err := topichandler.TDriver.PubPodUpdate(pods); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubPodUpdate error {%v}", err)
				log.Error(err)
				return
			}
			podlastupdatetime = time.Now().Add(time.Duration(md.info.PodUpdate.TimeOut) * time.Second)
		}

		time.Sleep(time.Millisecond * time.Duration(md.info.PodUpdate.DetectIntervalMill))
	}
}
