#!/bin/bash

echo "sleeping for 10 seconds"
sleep 10

mongo --host mongo-server.1:27017 << EOF
    load('/setup/init-replicaset.js')
EOF