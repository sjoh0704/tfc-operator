#!/bin/sh
# make generate manifests
img=192.168.9.12:5000/tfc-operator:$1
make docker-build IMG=$img
make docker-push IMG=$img 