module linklab/device-control-v2/edge-build

go 1.15

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/mholt/archiver/v3 v3.5.0
	github.com/sirupsen/logrus v1.8.1
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
