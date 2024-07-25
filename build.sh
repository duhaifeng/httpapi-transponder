#!/bin/bash
set -e

echo "building local os version ..."
echo ${PWD}
go build -o ./deploy/httpapi-transponder ./cmd/main.go
echo "build successfully ^o^"

echo "copy dependency config"
cd ./deploy
echo ${PWD}

cp  -r ../config ./
ls ./config
cp -r ../config_of_endpoint  ./
ls ./config_of_endpoint
echo "done"
