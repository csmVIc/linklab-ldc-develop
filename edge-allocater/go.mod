module linklab/device-control-v2/edge-allocater

go 1.15

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/nats-io/nats.go v1.11.0
	github.com/sirupsen/logrus v1.8.1
	k8s.io/api v0.18.19
	k8s.io/apimachinery v0.18.19
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
