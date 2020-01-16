#!/bin/sh

set -e

#用法提示
usage() {
    echo "Usage:"
    echo "  build.sh [-v --version]"
    echo "Description:"
    echo "    version 版本号."
    exit -1
}

if [[ $# -lt 2 ]]; then
    usage
fi

# 获取脚本执行时的选项
while getopts v: option
do
   case "${option}"  in
                v) Version=${OPTARG};;
                h) usage;;
                ?) usage;;
   esac
done

COMMID=$(git rev-parse --short HEAD)
NAME='cocoapods-cache-proxy-server'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main._version_=$Version -X main._commit_=$COMMID" -o ./docker/bin/$NAME-linux-amd64 && upx ./docker/bin/$NAME-linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main._version_=$Version -X main._commit_=$COMMID" -o ./docker/bin/$NAME-darwin-amd64 && upx ./docker/bin/$NAME-darwin-amd64