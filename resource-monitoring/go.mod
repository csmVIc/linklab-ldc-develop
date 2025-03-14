module linklab/device-control-v2/resource-monitoring

go 1.14

require (
	github.com/sirupsen/logrus v1.8.1
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/metrics v0.19.2
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
