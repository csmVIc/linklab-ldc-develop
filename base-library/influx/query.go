package influx

import (
	"context"
	"fmt"

	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// Query 查询数据
func (id *Driver) Query(query string) (*influxdb2api.QueryTableResult, error) {
	queryAPI := (*id.iclient).QueryAPI("")
	queryresult, err := queryAPI.Query(context.TODO(), query)
	if err != nil {
		err := fmt.Errorf("queryAPI.Query {%v} error {%v}", query, err)
		log.Error(err)
		return nil, err
	}

	return queryresult, nil
}
