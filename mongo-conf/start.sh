#!/bin/bash

# 1.清除上一阶段数据存储
sudo rm -rf data/mongo-server.1/*
sudo rm -rf data/mongo-server.2/*
sudo rm -rf data/mongo-server.3/*

# 2.mongodb rs集群配置
docker-compose -f docker-compose.yaml up -d
sleep 30
docker exec mongo-server.1 mongo /setup/init-mongo.js
docker rm mongo-conf_mongo-setup_1

# 3.mongodb 权限验证配置
docker-compose -f docker-compose-auth.yaml up -d