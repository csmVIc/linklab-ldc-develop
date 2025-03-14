# 命名空间
kubectl create namespace linklab

# 私有镜像仓库
kubectl create secret docker-registry linklab-aliyun --docker-server=registry.cn-hangzhou.aliyuncs.com --docker-username=qhgaoyi@gmail.com --docker-password=gaoy@aliyun --docker-email=yangg.china@outlook.com --namespace=linklab