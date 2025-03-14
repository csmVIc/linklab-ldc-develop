#!/bin/bash

# influx -username dataexporter -password 12 -precision rfc3339 -host 10.214.149.214 -port 30963 \
# -format json -database devicecontrollog -execute "select * from podmetrics where time > '$START_TIME' and time < '$END_TIME' group by container, podname tz('Asia/Shanghai')" \
# | jq . > $DATA_DIR/podmetrics.json