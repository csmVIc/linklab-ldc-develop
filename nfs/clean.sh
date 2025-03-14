#!/bin/bash

sudo rm -r data

sudo cp exports /etc/exports
sudo exportfs -a
sudo systemctl restart nfs-kernel-server