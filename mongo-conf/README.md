# LinkLab Device Control V2 Mongo 配置

## 启动

```c
chmod +x start.sh
./start.sh
```

## 调试

```c
// 进入到mongo的docker容器，执行如下命令
mongo -u DeviceControl -p 12 --authenticationDatabase linklab --host mongo-server-headless:27017
mongo -u Emqx -p 12 --authenticationDatabase linklab --host 127.0.0.1:27017
```

```shell
 db.mqtt_user.insert({"username": "emqx", "password": "7f384bb54f39fae699ef3a55d40dba2b87a5df2750148aa507585f0cb9d55380", "is_superuser": false, "salt": "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"})
```

## 问题

单节点情况下，执行事务时，会出现如下错误信息

```text
(IllegalOperation) Transaction numbers are only allowed on a replica set member or mongos
```

参考如下链接，确定问题和解决方法。在多副本集的情况下，mongo才会支持事务。

```text
1. https://www.youtube.com/watch?v=iq3ii2-hNi0
2. https://stackoverflow.com/questions/51238986/pymongo-transaction-errortransaction-numbers-are-only-allowed-on-a-replica-set
3. https://docs.mongodb.com/master/core/transactions/#transactions-and-replica-sets
```

解决方法如下

```shell
# 增加replSet参数
mongod --replSet rs
```

mongod启动之后，通过admin用户进入mongo命令行，执行如下命令

```shell
mongo -u LinkLab -p 12 --host 127.0.0.1:27017
>> rs.initiate()
...
{
    "operationTime" : Timestamp(1600095856, 7),
    "ok" : 0,
    "errmsg" : "already initialized",
    "code" : 23,
    "codeName" : "AlreadyInitialized",
    "$clusterTime" : {
        "clusterTime" : Timestamp(1600095856, 7),
        "signature" : {
            "hash" : BinData(0,"O1lpBnweEByPuEb9cpbLjdWBlSY="),
            "keyId" : NumberLong("6872359371985125380")
        }
    }
}
...
>> rs.conf()
...
{
	"_id" : "rs0",
	"version" : 1,
	"protocolVersion" : NumberLong(1),
	"writeConcernMajorityJournalDefault" : true,
...
```

最后mongo客户端代码的uri应该添加replicaSet参数

```go
    uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", ci.User, ci.Password, ci.Host, ci.Port, ci.Db)
    if ci.Single {
        uri = fmt.Sprintf("%s?replicaSet=rs", uri)
    }
```

helm repo add bitnami https://charts.bitnami.com/bitnami