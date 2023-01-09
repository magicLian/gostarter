#!/usr/bin/env bash

echo "start gen swagger"
cd pkg/api

swag fmt -d ./ --exclude ./api.go

swag init --pd -g api.go -o ../docs

