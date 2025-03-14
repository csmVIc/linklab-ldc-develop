module linklab/device-control-v2/riot-worker

go 1.14

require (
	github.com/mholt/archiver/v3 v3.5.0
	github.com/nats-io/nats.go v1.11.0
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.5.3
	linklab/device-control-v2/base-library v0.0.0-00010101000000-000000000000
)

replace linklab/device-control-v2/base-library => ../base-library
