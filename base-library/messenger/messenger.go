package messenger

import (
	"errors"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

// Driver 消息队列客户端
type Driver struct {
	info     *MInfo
	natsConn *nats.Conn
	stanConn stan.Conn
	closed   bool
}

// Driver 消息队列客户端全局实例
var (
	Mdriver *Driver
)

// GetNatsConn 返回连接实例
func (md *Driver) GetNatsConn() (*nats.Conn, error) {

	if md.natsConn == nil {
		return nil, errors.New("nats conn nil err")
	}

	if md.natsConn.IsConnected() == false {
		if err := md.init(); err != nil {
			return nil, err
		}
	}

	return md.natsConn, nil
}

// GetStanConn 返回连接实例
func (md *Driver) GetStanConn() (*stan.Conn, error) {

	if md.natsConn == nil || md.stanConn == nil {
		return nil, errors.New("nats/stan conn nil err")
	}

	if md.natsConn.IsConnected() == false {
		if err := md.init(); err != nil {
			log.Errorf("md.init() error {%v}", err)
			return nil, err
		}
	}

	return &md.stanConn, nil
}

func (md *Driver) init() error {

	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("init os.Hostname error {%v}", err)
		return err
	}
	pid := os.Getpid()
	connname := fmt.Sprintf("%s-%v", hostname, pid)

	md.natsConn, err = nats.Connect(md.info.Client.URL, nats.Name(connname),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("nats client disconnected {%v}", err)
		}), nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("nats client reconnected")
		}), nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("client closed")
			md.closed = true
		}))
	if err != nil {
		log.Errorf("init nats.Connect error {%v}", err)
		return err
	}

	if md.info.Client.NeedStan {
		md.stanConn, err = stan.Connect(md.info.Client.ClusterID, connname, stan.NatsConn(md.natsConn))
		if err != nil {
			log.Errorf("init stan.Connect error {%v}", err)
			return err
		}
	}

	if isOk := md.natsConn.IsConnected(); isOk == false {
		err := errors.New("nats conn error")
		log.Error(err)
		return err
	}

	md.closed = false
	return nil
}

// New 创建消息队列客户端实例
func New(i *MInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("info is nil error")
	}
	md := &Driver{info: i, natsConn: nil, stanConn: nil, closed: true}
	if err := md.init(); err != nil {
		log.Errorf("init error {%v}", err)
		return nil, err
	}
	return md, nil
}
