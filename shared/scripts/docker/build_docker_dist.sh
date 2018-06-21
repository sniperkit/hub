#!/bin/sh
set -x
set -e

env GOOS=linux go build -v ./shared/docker/alpine/dist/hub-linux ./cmd/hub/*.go
docker build --force-rm -t sniperkit/hub.cli:dist-3.7-alpine --no-cache -f ./shared/docker/alpine/dist/dockerfile-alpine3.7