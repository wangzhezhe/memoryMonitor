FROM golang:1.7-alpine
MAINTAINER zhewang@daocloud.io

COPY ./memoryMonitor /usr/bin/memoryMonitor

ENTRYPOINT ["/usr/bin/memoryMonitor"]
 
