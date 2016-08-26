#!/bin/bash
export GOPATH=$HOME/work
wget -N http://updates.jenkins-ci.org/update-center.json
cat update-center.json | sed 1d | sed '$ d'>json
go run jenkins_plugin_mirror.go
