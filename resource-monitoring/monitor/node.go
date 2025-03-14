package monitor

import (
	"context"
	"linklab/device-control-v2/base-library/logger"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// node监控
func (md *Driver) nodemonitor() {
	defer md.exitwg.Done()

	for {
		nodemetricslist := &v1beta1.NodeMetricsList{}
		err := md.clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Into(nodemetricslist)
		if err != nil {
			log.Errorf("monitor.Driver.nodemonitor nodemetricslist get error {%v}", err)
			return
		}

		nodelist, err := md.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("monitor.Driver.nodemonitor nodelist get error {%v}", err)
			return
		}

		// node -> allocatable cpu memory
		type resource struct {
			cpu    int64
			memory int64
			arch   string
		}
		nmap := make(map[string]resource)
		for _, n := range nodelist.Items {
			nmap[n.GetName()] = resource{
				cpu:    n.Status.Allocatable.Cpu().MilliValue(),
				memory: n.Status.Allocatable.Memory().MilliValue(),
				arch:   n.Labels["kubernetes.io/arch"],
			}
		}

		for _, nm := range nodemetricslist.Items {

			// 数据库日志记录
			tags := map[string]string{
				"nodename": nm.GetName(),
			}
			fields := map[string]interface{}{
				"record.time":        time.Now().UnixNano(),
				"cpu.usage":          nm.Usage.Cpu().MilliValue(),
				"memory.usage":       nm.Usage.Memory().MilliValue(),
				"cpu.allocatable":    nmap[nm.GetName()].cpu,
				"memory.allocatable": nmap[nm.GetName()].memory,
				"arch":               nmap[nm.GetName()].arch,
			}

			if err := logger.Ldriver.WriteLog("nodemetrics", tags, fields); err != nil {
				log.Errorf("database log {%v} error", err)
			}
		}

		time.Sleep(time.Duration(md.info.MetricsMap["nodemetrics"].Interval) * time.Second)
	}
}
