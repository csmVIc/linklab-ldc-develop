module linklab/device-control-v2/user-service

go 1.14

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-resty/resty/v2 v2.6.0
	github.com/gorilla/websocket v1.4.2
	github.com/nats-io/nats.go v1.11.0
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.5.4
	linklab/device-control-v2/base-library v0.0.1
)

replace linklab/device-control-v2/base-library => ../base-library
