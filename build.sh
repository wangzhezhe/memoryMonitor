#!/bin/bash

set -ex
#check the go command manually

CGO_ENABLED=0 go install -a github.com/memoryMonitor

cp $GOPATH/bin/memoryMonitor .

sudo docker build -t daocloud.io/cform_monitor:v0.1 .

sudo docker run -id -v /sys:/sys: -v /proc:/proc daocloud.io/cform_monitor:v0.1 -nodeip=<nodeip>