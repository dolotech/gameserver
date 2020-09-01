#!/bin/bash
#compile json 信息
# usage ./build-linux.sh 1.0.13
curDir=`pwd`

cd   $curDir/..

# export GOPATH=`pwd`
# export GOARCH=amd64
# export GOOS=linux


cd bin
go run  ../src/main.go --alsologtostderr --log_dir=./logs


read -p "Press any key to continue." var
