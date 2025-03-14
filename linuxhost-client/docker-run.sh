#!/bin/bash

sudo rm -rf tmp
mkdir tmp

sudo rm -rf workspace
mkdir workspace

docker-compose -f docker-compose.yaml up -d
