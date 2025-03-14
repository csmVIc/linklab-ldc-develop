#!/bin/bash

sudo rm -rf tmp
mkdir tmp

docker-compose -f docker-compose.yaml up -d
