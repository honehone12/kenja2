#!/bin/bash

pwd
go get -u ./...
cd bin
pwd
go get -u ./...
cd ../lambda
pwd
go get -u ./...
cd ..
