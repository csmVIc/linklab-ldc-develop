package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) resource() error {

	defer func() {
		// 如果该函数退出,则将结束运行标志位设置为true
		md.endrun = true
		log.Error("resource func set endrun")
	}()

	for md.endrun == false {

		// 获取边缘设备资源占用
		edgeNodeResource, err := edgenode.EDriver.GetEdgeNodeResourceList()
		if err != nil {
			err = fmt.Errorf("edgenode.EDriver.GetEdgeNodeResourceList error {%v}", err)
			log.Error(err)
			return err
		}
		if edgeNodeResource != nil && len(edgeNodeResource.EdgeNodes) > 0 {
			// 上传边缘设备资源占用
			if err := topichandler.TDriver.PubEdgeNodeResource(edgeNodeResource); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubEdgeNodeResource error {%v}", err)
				log.Error(err)
				return err
			}
		}

		// 获取Pod资源占用
		podResource, err := edgenode.EDriver.GetPodResourceList()
		if err != nil {
			err = fmt.Errorf("edgenode.EDriver.GetPodResourceList error {%v}", err)
			log.Error(err)
			return err
		}
		if podResource != nil && len(podResource.Pods) > 0 {
			// 上传Pod资源占用
			if err := topichandler.TDriver.PubPodResource(podResource); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubPodResource error {%v}", err)
				log.Error(err)
				return err
			}
		}

		time.Sleep(time.Second * time.Duration(md.info.ResourseUpdate.Interval))
	}

	return nil
}
