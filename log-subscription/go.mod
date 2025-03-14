module linklab/device-control-v2/log-subscription

go 1.15

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.5.4
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
