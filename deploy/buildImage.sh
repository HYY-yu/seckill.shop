#!/bin/bash

#
# This script builds a Docker image for the controller
# and pushes it to the local Docker repository
# It assumes that the following environment variables are set
# TRAVIS_TAG - if the current build is for a tag push, this should be the tag name, otherwise
#              this is assumed to be empty
# TRAVIS_BUILD_DIR - the absolute path name of the cloned repository
# The script will create a file  $TRAVIS_BUILD_DIR/tag which contains the tag
# used for the push

# Fail if a line fails
set -ev

#
# Get short form of git hash for current commit
#
hash=$(git log --pretty=format:'%h' -n 1)
serviceName='shop'
dockerName="registry.cn-hangzhou.aliyuncs.com/hyy_yu/seckill.${serviceName}"

#
# Determine tag. If the build is from a tag push, use tag name, otherwise
# use commit hash
#
if [ "X$TRAVIS_TAG" == "X" ]; then
  tag=$hash
else
  tag=$TRAVIS_TAG
fi


#
# Create image locally
cd $TRAVIS_BUILD_DIR/internal/service/$serviceName
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o sk_$serviceName ./cmd/main.go
docker build --rm -f $TRAVIS_BUILD_DIR/deploy/Dockerfile --build-arg serviceName=$serviceName -t $dockerName:$tag .
docker tag $dockerName:$tag $dockerName:latest

#
# Return tag
#
echo $tag > $TRAVIS_BUILD_DIR/tag