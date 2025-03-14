# LinkLab Device Control V2 Influxdb 配置

## 调试

```shell
// 进入到influxdb的docker容器，执行如下命令
influx -username linklab -password 12 -precision rfc3339
```

delete from burnresult, clientconnect, clientdisconnect, endrun, entertaskswait, filedownload, fileupload, nodemetrics, podmetrics, taskallocate, userlogin, burntask, compilelog, nodemetrics, podmetrics