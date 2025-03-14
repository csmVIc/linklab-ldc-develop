module linklab/device-control-v2/decision-maker

go 1.14

require (
	github.com/go-redis/redis/v8 v8.10.0
	github.com/nats-io/nats.go v1.11.0
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
