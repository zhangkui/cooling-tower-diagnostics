#!/usr/bin/env sh
set -eu
name=${1:?image name required}
platform=${2:-linux/amd64}
docker build --platform "$platform" -f benzhi.Dockerfile -t "benzhi/$name:latest" .
