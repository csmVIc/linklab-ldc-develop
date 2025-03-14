# LinkLab Device Control V2 Redis 配置

## 版本选择

```shell
docker pull redis:6.0.7-alpine3.12
```

## 运行

### docker

### kubernetes

redis-cli -h redis-server -p 6379 -a 41eb37269c0525a25ea19b59f65d12414f103556e5ab1f7a7e7d3e3553ee9941

SENTINEL get-master-addr-by-name master

redis-cli -h redis-server -p 26379 -a 41eb37269c0525a25ea19b59f65d12414f103556e5ab1f7a7e7d3e3553ee9941

redis-cli -h redis-server-headless -p 26379 -a 41eb37269c0525a25ea19b59f65d12414f103556e5ab1f7a7e7d3e3553ee9941

kubectl run --namespace linklab redis-server-client --rm --tty -i --restart='Never' \
    --env REDIS_PASSWORD=$REDIS_PASSWORD \
   --image docker.io/bitnami/redis:6.0.7 -- bash

redis-cli -h redis-server -p 26379 -a $REDIS_PASSWORD

redis-cli -h redis-server-headless -a $REDIS_PASSWORD
   redis-cli -h redis-server-slave -a $REDIS_PASSWORD
