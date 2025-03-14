#!/bin/bash

docker-compose -f docker-compose.yaml down

sudo rm -rf tmp
mkdir tmp