#!/bin/sh

set -e

echo "Create Scheduler..."
gcloud beta scheduler jobs create app-engine example-job \
  --schedule "0 15 * * *" \
  --description "Example job" \
  --service "tasks" \
  --relative-url "/cron" \
  --time-zone "Asia/Tokyo" \
  --max-retry-attempts 5

echo "Create Tasks Queue..."
gcloud beta tasks queues create-app-engine-queue example-queue --max-attempts 5

