module linklab/device-control-v2/xiuos-worker

go 1.14

require (
	github.com/mholt/archiver/v3 v3.5.0
	github.com/nats-io/nats.go v1.10.1-0.20210228004050-ed743748acac
	github.com/sirupsen/logrus v1.7.0
	go.mongodb.org/mongo-driver v1.4.4
	linklab/device-control-v2/base-library v0.0.0-00010101000000-000000000000
)

replace linklab/device-control-v2/base-library => ../base-library
