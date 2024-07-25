#!/bin/bash

APP_NAME=httpapi-transponder

git pull
\rm -rf $APP_NAME
go build -o $APP_NAME \
    -ldflags "-X main.BuildCommitID=`git rev-parse HEAD` -X main.BuildTime=`date '+%Y-%m-%d'`  -X main.BuildBranch=`git rev-parse --abbrev-ref HEAD`" \
    cmd/server_starter.go
pidof $APP_NAME | xargs kill
nohup ./$APP_NAME  &

pidof %APP_NAME
tailf nohup.out
