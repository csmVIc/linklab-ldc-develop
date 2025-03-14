module linklab/device-control-v2/jlink-client

go 1.15

require (
	github.com/albenik/go-serial/v2 v2.5.0
	github.com/sirupsen/logrus v1.8.1
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
