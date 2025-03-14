#!/bin/bash
IS_HA=false

if [ $IS_HA == true ] 
then
    # 高可用集群配置 
    sudo kubeadm init --upload-certs --config kubeadm-config.yaml
else
    # 单控制平面节点配置
    sudo kubeadm init --config kubeadm-config.yaml
fi

# kubectl配置文件权限
rm -r $HOME/.kube
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 允许控制平面节点运行任务Pod
kubectl taint nodes --all node-role.kubernetes.io/master-

# 配置weave网络
# kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
kubectl apply -f kube-flannel.yml
kubectl apply -f ingress.yaml
# 资源监控
kubectl apply -f components.yaml

# 命名空间
kubectl create namespace linklab

# 私有镜像仓库
kubectl create secret docker-registry linklab-aliyun --docker-server=registry.cn-hangzhou.aliyuncs.com --docker-username=qhgaoyi@gmail.com --docker-password=gaoy@aliyun --docker-email=yangg.china@outlook.com --namespace=linklab

# 节点资源监控间隔
# sudo cp kubeadm-flags.env /var/lib/kubelet/kubeadm-flags.env
# sudo systemctl daemon-reload && sudo systemctl restart kubelet

# nginx ingress 配置
# kubectl apply -f ingress.yaml

# helm配置
# sudo cp helm /usr/local/bin/helm
# helm repo add bitnami https://charts.bitnami.com/bitnami
# helm repo add stable https://charts.helm.sh/stable     
# helm repo add emqx https://repos.emqx.io/charts


