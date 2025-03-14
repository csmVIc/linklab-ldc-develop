module linklab/device-control-v2/edge-client

go 1.14

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/gorilla/websocket v1.4.2
	github.com/sirupsen/logrus v1.8.1
	k8s.io/api v0.18.19
	k8s.io/apimachinery v0.18.19
	k8s.io/client-go v0.18.19
	k8s.io/metrics v0.18.19
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
