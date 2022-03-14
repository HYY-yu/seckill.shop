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

serviceName='shop'
dockerName="registry.cn-hangzhou.aliyuncs.com/hyy_yu/seckill.${serviceName}"


#
# Login to Docker hub
#
#echo "$TENCENT_DOCKER_PASSWORD" | docker login ccr.ccs.tencentyun.com --username $TENCENT_DOCKER_USER --password-stdin
echo "$TENCENT_DOCKER_PASSWORD" | docker login --username=13718640725 registry.cn-hangzhou.aliyuncs.com --password-stdin

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
cd /tmp
git config --global user.name feng.yu
git config --global user.email 690174435@qq.com
git clone https://$GITHUB_USER:$GITHUB_PASSWORD@github.com/HYY-yu/seckill.$serviceName.chart.git
cd seckill.$serviceName.chart

#
# Get current version from Chart file and remove build tag
#
current_version=$(cat Chart.yaml | grep "version" | awk '{ print $2 }' | sed 's/-dev[a-z,0-9]*//')
echo "Current chart version: $current_version"
#
# Update chart version and appVersion  in chart file
#
if [ "X$TRAVIS_TAG" != "X" ]; then
  chart_version=$TRAVIS_TAG
else
  chart_version="$current_version-dev$tag"
fi
echo "Using chart version $chart_version"
cat Chart.yaml | sed "s/version.*/version: $chart_version/"  | sed "s/appVersion.*/appVersion: $tag/" > /tmp/Chart.yaml.patched
cp /tmp/Chart.yaml.patched Chart.yaml


git add --all
git config remote.origin.url https://$GITHUB_USER:$GITHUB_PASSWORD@github.com/HYY-yu/seckill.$serviceName.chart
git commit -m "Automated deployment of chart version $chart_version"
git push origin main