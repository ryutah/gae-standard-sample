#!/bin/bash

echo "Start Server With: $@"

dev_appserver.py \
  --host 0.0.0.0 \
  --admin_host 0.0.0.0 \
  $@
