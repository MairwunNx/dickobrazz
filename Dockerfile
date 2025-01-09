# syntax=docker/dockerfile:1

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.23.4 AS builder
COPY go.mod go.sum ./
RUN go mod download -x
COPY *.go ./
COPY *.env ./
COPY .version ./
COPY .git ./
COPY application/*.go ./application/
COPY application/geo/*.go ./application/geo/
ENV GOCACHE=/root/.cache/go-build
RUN go build -o /app/dickobot -a -installsuffix cgo -gcflags "all=-N -l" -tags timetzdata -ldflags="-s -w -X dickobot/application.Version=$(cat .version) -X dickobot/application.BuildAt=$(date +%Y-%m-%d_%H:%M:%S) -X dickobot/application.GoVersion=$(go version | awk '{print $3}') -X dickobot/application.BuildRv=$(git describe --always --long)"
RUN chmod +x /app/dickobot

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/dickobot dickobot
CMD ["./dickobot"]