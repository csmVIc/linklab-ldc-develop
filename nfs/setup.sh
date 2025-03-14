#!/bin/bash

sudo apt install -y nfs-kernel-server
sudo apt install -y nfs-common

sudo mkdir -p data
sudo chown nobody:nogroup data
sudo chmod 777 data

mkdir data/redis
mkdir data/redis/data-0
mkdir data/redis/data-1
mkdir data/redis/data-2
sudo chown -R 1001:1001 data/redis/data-0
sudo chown -R 1001:1001 data/redis/data-1
sudo chown -R 1001:1001 data/redis/data-2

mkdir data/influxdb
mkdir data/influxdb/data

mkdir data/stan
mkdir data/stan/data-0
mkdir data/stan/data-1
mkdir data/stan/data-2

mkdir data/mongodb
mkdir data/mongodb/data-0
mkdir data/mongodb/data-1
mkdir data/mongodb/data-2

mkdir data/linuxhost-client
mkdir data/linuxhost-client/writespace-0
mkdir data/linuxhost-client/writespace-0/tmp
mkdir data/linuxhost-client/writespace-0/workspace
cp -r data/linuxhost-client/writespace-0 data/linuxhost-client/writespace-1

mkdir data/log
mkdir data/log/data
mkdir data/log/data/login-authentication
mkdir data/log/data/file-cache
mkdir data/log/data/client-monitoring
mkdir data/log/data/decision-maker
mkdir data/log/data/device-service
mkdir data/log/data/user-service
mkdir data/log/data/resource-monitoring
mkdir data/log/data/device-monitoring
mkdir data/log/data/log-subscription
mkdir data/log/data/linuxhost-client

mkdir data/compile-worker
mkdir data/compile-worker/alios-haas-worker
mkdir data/compile-worker/alios-haas-worker/init
mkdir data/compile-worker/alios-haas-worker/writespace.0
mkdir data/compile-worker/alios-haas-worker/writespace.1

mkdir data/compile-worker/arduino-worker
mkdir data/compile-worker/arduino-worker/init
mkdir data/compile-worker/arduino-worker/writespace.0
mkdir data/compile-worker/arduino-worker/writespace.1

mkdir data/compile-worker/tinysim-worker
mkdir data/compile-worker/tinysim-worker/init
mkdir data/compile-worker/tinysim-worker/writespace.0
mkdir data/compile-worker/tinysim-worker/writespace.1

mkdir data/compile-worker/alios-worker
mkdir data/compile-worker/alios-worker/init
mkdir data/compile-worker/alios-worker/writespace.0
mkdir data/compile-worker/alios-worker/writespace.1

mkdir data/compile-worker/contiki-ng-worker
mkdir data/compile-worker/contiki-ng-worker/init
mkdir data/compile-worker/contiki-ng-worker/writespace.0
mkdir data/compile-worker/contiki-ng-worker/writespace.1

mkdir data/compile-worker/mbed-worker
mkdir data/compile-worker/mbed-worker/init
mkdir data/compile-worker/mbed-worker/writespace.0
mkdir data/compile-worker/mbed-worker/writespace.1

mkdir data/compile-worker/esp8266duino-worker
mkdir data/compile-worker/esp8266duino-worker/init
mkdir data/compile-worker/esp8266duino-worker/writespace.0
mkdir data/compile-worker/esp8266duino-worker/writespace.1

mkdir data/compile-worker/stm32duino-worker
mkdir data/compile-worker/stm32duino-worker/init
mkdir data/compile-worker/stm32duino-worker/writespace.0
mkdir data/compile-worker/stm32duino-worker/writespace.1

mkdir data/compile-worker/esp32duino-worker
mkdir data/compile-worker/esp32duino-worker/init
mkdir data/compile-worker/esp32duino-worker/writespace.0
mkdir data/compile-worker/esp32duino-worker/writespace.1

mkdir data/compile-worker/esp32duino-virtual-worker
mkdir data/compile-worker/esp32duino-virtual-worker/init
mkdir data/compile-worker/esp32duino-virtual-worker/writespace.0
mkdir data/compile-worker/esp32duino-virtual-worker/writespace.1

mkdir data/compile-worker/stm32-gcc-worker
mkdir data/compile-worker/stm32-gcc-worker/init
mkdir data/compile-worker/stm32-gcc-worker/writespace.0
mkdir data/compile-worker/stm32-gcc-worker/writespace.1

mkdir data/compile-worker/harmonyos-worker
mkdir data/compile-worker/harmonyos-worker/init
mkdir data/compile-worker/harmonyos-worker/writespace.0
mkdir data/compile-worker/harmonyos-worker/writespace.1

mkdir data/compile-worker/riot-worker
mkdir data/compile-worker/riot-worker/init
mkdir data/compile-worker/riot-worker/writespace.0
mkdir data/compile-worker/riot-worker/writespace.1

mkdir data/compile-worker/alios-edu-worker
mkdir data/compile-worker/alios-edu-worker/init
mkdir data/compile-worker/alios-edu-worker/writespace.0
mkdir data/compile-worker/alios-edu-worker/writespace.1

mkdir data/compile-worker/openharmony-worker
mkdir data/compile-worker/openharmony-worker/init
mkdir data/compile-worker/openharmony-worker/writespace.0
mkdir data/compile-worker/openharmony-worker/writespace.1

mkdir data/compile-worker/xiuos-worker
mkdir data/compile-worker/xiuos-worker/init
mkdir data/compile-worker/xiuos-worker/writespace.0
mkdir data/compile-worker/xiuos-worker/writespace.1

mkdir data/static-file-service
mkdir data/static-file-service/data

echo "$(pwd)/data *(rw,sync,no_subtree_check,no_root_squash)" | sudo tee -a /etc/exports
sudo exportfs -a
sudo systemctl restart nfs-kernel-server