FROM gcr.io/cloud-builders/gcloud-slim

RUN add-apt-repository -y ppa:longsleep/golang-backports \
 && apt-get update \
 && apt-get install -y golang-1.12

RUN gcloud --quiet components install app-engine-go

RUN mkdir /go

COPY ./serve.sh /local/bin/serve.sh

ENV PATH=/go/bin:/usr/lib/go-1.12/bin/:$PATH \
    GOPATH=/go

ENTRYPOINT ["/local/bin/serve.sh"]
