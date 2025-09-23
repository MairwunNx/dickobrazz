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
COPY application/database/*.go ./application/database/
COPY application/datetime/*.go ./application/datetime/
COPY application/geo/*.go ./application/geo/
COPY application/logging/*.go ./application/logging/
COPY application/timings/*.go ./application/timings/
RUN go build -o /app/dickobot -a -installsuffix cgo -gcflags "all=-N -l" -tags timetzdata -ldflags="-s -w -X dickobot/application/logging.Version=$(cat .version) -X dickobot/application/logging.BuildAt=$(date +%Y-%m-%d_%H:%M:%S) -X dickobot/application/logging.GoVersion=$(go version | awk '{print $3}') -X dickobot/application/logging.BuildRv=$(git describe --always --long)"
RUN chmod +x /app/dickobot

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/dickobot dickobot
CMD ["./dickobot"]