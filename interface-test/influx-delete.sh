#!/bin/bash

array=("burnresult" "clientconnect" "clientdisconnect" "devicelist" "endrun" "entertaskswait" "filedownload" "fileupload" "nodemetrics" "podmetrics" "taskallocate" "userlogin")

for elem in ${array[@]}
do
  echo $elem
  influx -username linklab -password 12 -precision rfc3339 -host 10.214.149.214 -port 30963 \
  -format json -database devicecontrollog -execute "delete from $elem"
done

