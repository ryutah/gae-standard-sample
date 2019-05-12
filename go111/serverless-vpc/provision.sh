#!/bin/sh

set -eux

gcloud services enable vpcaccess.googleapis.com

gcloud beta compute networks vpc-access connectors create example \
  --network default \
  --region us-central1 \
  --range 10.8.0.0/28

gcloud deployment-manager deployments create \
  --automatic-rollback-on-error \
  --template resource.jinja \
  serverless-vpc-examlple
