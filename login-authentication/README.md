接口

用户登录

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"UserTest","password":"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}' \
  http://localhost:8082/user/login
```

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"emqx","password":"efa1f375d76194fa51a3556a97e641e61685f914d446979da50a551a4333ffd7"}' \
  http://localhost:8082/client/login
```

```shell
curl --header "Content-Type: application/json" \
  --header "Authorization: e3ab1c5c1368c432b232cbc2a071fc431faa85e9c78f3ad31bb97f562d82e64d" \
  --request POST \
  http://localhost:8082/user/signout
```


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"UserTest","password":"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918"}' \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/login-authentication/user/login