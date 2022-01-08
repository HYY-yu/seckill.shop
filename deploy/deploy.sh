set -ev

dockerName='ccr.ccs.tencentyun.com/seckill_srv/seckill-shop'

#
# Login to Docker hub
#
echo "$TENCENT_DOCKER_PASSWORD" | docker login --username $TENCENT_DOCKER_USER --password-stdin

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
