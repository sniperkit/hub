#!/bin/sh
set -x
set -e

env GOOS=linux go build -v ./docker/hub-dev
docker build --force-rm -t sniperkit/hub.cli:dev-3.7-alpine --no-cache -f ./docker/alpine/dev/dockerfile-alpine3.7