#!/bin/bash
sudo kubeadm reset
sudo iptables -F && sudo iptables -t nat -F && sudo iptables -t mangle -F && sudo iptables -X
sudo ipvsadm -C
sudo iptables -P INPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -t nat -F
sudo iptables -t mangle -F
sudo iptables -F
sudo iptables -X
sudo rm /opt/cni/bin/weave-*
sudo ./weave reset --force
rm -r $HOME/.kube
sudo rm -rf /etc/cni/net.d/*
sudo rm -rf /var/lib/cni/
sudo ip link delete cni0
sudo ip link delete flannel.1
sudo systemctl daemon-reload
sudo systemctl restart docker