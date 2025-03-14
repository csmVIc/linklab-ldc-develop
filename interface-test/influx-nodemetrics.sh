#!/bin/bash

# influx -username dataexporter -password 12 -precision rfc3339 -host 10.214.149.214 -port 30963 \
# -format json -database devicecontrollog -execute "select * from nodemetrics where time > '$START_TIME' and time < '$END_TIME' group by nodename tz('Asia/Shanghai')" \
# | jq . > $DATA_DIR/nodemetrics.json