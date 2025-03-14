#!/bin/sh
ip=$1
filepath=$2
filename=$3

tftp ${ip}  << !
binary
put ${filepath} ${ip}:/data/${filename}
quit