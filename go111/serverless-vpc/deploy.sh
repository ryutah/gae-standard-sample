#!/bin/sh

cd `dirname $0`
wd=`pwd`

set -eux

cd ./k8s-service
make apply

cd ${wd}
cd ./appengine-service
if [[ -f ./app.yaml ]]; then
  rm ./app.yaml
fi
sed 's/\$GOOGLE_CLOUD_PROJECT'"/${GOOGLE_CLOUD_PROJECT}/" app.yaml.template > app.yaml
gcloud --quiet beta app deploy .
