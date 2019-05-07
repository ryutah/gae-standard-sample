FROM gcr.io/cloud-builders/gcloud-slim

RUN apt-get update \
 && apt-get install -y openjdk-9-jre

RUN gcloud --quiet components install beta cloud-firestore-emulator

COPY ./start.sh /local/bin/start.sh

ENV FIRESTORE_PORT=5000

ENTRYPOINT /local/bin/start.sh
