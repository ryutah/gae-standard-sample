FROM golang:1.12-alpine as builder
COPY . /app
WORKDIR /app
RUN go build -o app .

FROM alpine
COPY --from=builder /app/app /usr/local/bin/
CMD ["/usr/local/bin/app"]
