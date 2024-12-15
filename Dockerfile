# syntax=docker/dockerfile:1

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.23.4 AS builder
COPY go.mod go.sum ./
ENV GOPATH=""
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /app/dickobot -a -installsuffix cgo -ldflags="-s -w -X application.Version="$(cat .version)" -X application.BuildAt="$(date +%Y-%m-%d_%H:%M:%S)" -X application.BuildRv="$(git describe --always --long)""
RUN chmod +x /app/dickobot

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/dickobot dickobot
CMD ["/dickobot"]