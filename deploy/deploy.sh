set -e
#
# Push images to Docker hub and update Helm chart repository
#
# This script assumes that there is a file tag
# which contains the name of the tag to use for the push
#
# We assume that the following environment variables are set
# DOCKER_USER               User for docker hub
# DOCKER_PASSWORD           Password
# TRAVIS_BUILD_DIR          Travis build directory
# TRAVIS_TAG                Tag if the build is caused by a git tag

dockerName='ccr.ccs.tencentyun.com/seckill_srv/seckill-shop'

#
# Login to Docker hub
#
echo "$TENCENT_DOCKER_PASSWORD" | docker login ccr.ccs.tencentyun.com --username $TENCENT_DOCKER_USER --password-stdin

#
# Get tag
#
tag=$(cat $TRAVIS_BUILD_DIR/tag)

#
# Push images
#
docker push $dockerName:$tag
docker push $dockerName:latest

#
# Now clone into the repository that contains the Helm chart
#
