module linklab/device-control-v2/interface-test

go 1.15

require (
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/mholt/archiver/v3 v3.3.0
	github.com/sirupsen/logrus v1.6.0
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
