# syntax=docker/dockerfile:1

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.23.4 AS builder

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY *.go ./
COPY *.env ./
COPY .version ./
COPY .git ./
COPY application/*.go ./application/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /app/dickobot -a -installsuffix cgo -gcflags "all=-N -l" \
    -tags timetzdata -ldflags="-s -w -X dickobot/application.Version=$(cat .version) -X dickobot/application.BuildAt=$(date +%Y-%m-%d_%H:%M:%S) -X dickobot/application.GoVersion=$(go version | awk '{print $3}') -X dickobot/application.BuildRv=$(git describe --always --long)"

RUN chmod +x /app/dickobot

FROM busybox
EXPOSE 10000
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/dickobot dickobot
COPY --from=builder /go/bin/dlv /

CMD ["./dickobot"]