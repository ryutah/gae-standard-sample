#!/bin/sh

gcloud beta emulators firestore start --host-port 0.0.0.0:${FIRESTORE_PORT}
