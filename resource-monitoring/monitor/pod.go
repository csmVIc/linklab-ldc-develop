package monitor

import (
	"context"
	"time"

	"linklab/device-control-v2/base-library/logger"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// pod监控
func (md *Driver) podmonitor() {
	defer md.exitwg.Done()

	for {
		podmetricslist := &v1beta1.PodMetricsList{}
		err := md.clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/" + md.info.MetricsMap["podmetrics"].Namespace + "/pods").Do(context.TODO()).Into(podmetricslist)
		if err != nil {
			log.Errorf("podmetricslist get error {%v}", err)
			return
		}

		podlist, err := md.clientset.CoreV1().Pods(md.info.MetricsMap["podmetrics"].Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("podlist get error {%v}", err)
			return
		}

		// namespace/podname -> node
		pmap := make(map[string]string)
		for _, p := range podlist.Items {
			pmap[p.GetNamespace()+"/"+p.GetName()] = p.Spec.NodeName
		}

		for _, pm := range podmetricslist.Items {
			for _, c := range pm.Containers {
				// 数据库日志记录
				tags := map[string]string{
					"podname":   pm.GetName(),
					"namespace": pm.GetNamespace(),
					"container": c.Name,
				}
				fields := map[string]interface{}{
					"record.time": time.Now().UnixNano(),
					"node":        pmap[pm.GetNamespace()+"/"+pm.GetName()],
					"cpu":         c.Usage.Cpu().MilliValue(),
					"memory":      c.Usage.Memory().MilliValue(),
				}

				if err := logger.Ldriver.WriteLog("podmetrics", tags, fields); err != nil {
					log.Errorf("database log {%v} error", err)
				}
			}
		}

		time.Sleep(time.Duration(md.info.MetricsMap["podmetrics"].Interval) * time.Second)
	}
}
