package influxv1

import (
	"fmt"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

// Query 查询数据
func (id *Driver) Query(query string) (*client.Response, error) {
	q := client.Query{
		Command:  query,
		Database: id.info.Client.DataBase,
	}
	queryRes, err := (*id.iclient).Query(q)
	if err != nil {
		err = fmt.Errorf("(*id.iclient).Query error {%v}", err)
		log.Error(err)
		return nil, err
	}
	if queryRes.Error() != nil {
		err = fmt.Errorf("queryRes.Error {%v}", queryRes.Error())
		log.Error(err)
		return nil, err
	}
	return queryRes, nil
}
