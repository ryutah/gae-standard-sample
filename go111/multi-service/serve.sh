#!/bin/sh

cd `dirname $0`
wd=`pwd`

set -eux


cd ${wd}/service1
if [ -f cron.yaml ]; then
  unlink cron.yaml
fi
ln -s ../cron.yaml ./

cd ${wd}
dev_appserver.py ./service1 ./service2 ./dispatch.yaml
