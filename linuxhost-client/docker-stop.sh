#!/bin/bash

docker-compose -f docker-compose.yaml down

sudo rm -rf tmp
mkdir tmp

sudo rm -rf workspace
mkdir workspace